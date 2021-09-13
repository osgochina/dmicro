package process

import (
	"fmt"
	loggerv2 "github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/supervisor/logger"
	"github.com/osgochina/dmicro/supervisor/procconf"
	"github.com/osgochina/dmicro/supervisor/signals"
	"io"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Process struct {
	procEntry *procconf.ProcEntry
	cmd       *exec.Cmd
	startTime time.Time
	stopTime  time.Time
	// 进程的当前状态
	state State
	// 正在启动的时候，该值为true
	inStart bool
	// 用户主动关闭的时候，该值为true
	stopByUser bool
	//启动的次数
	retryTimes *int32
	lock       sync.RWMutex
	stdin      io.WriteCloser
	StdoutLog  logger.Logger
	StderrLog  logger.Logger
}

// NewProcess 创建进程对象
func NewProcess(procEntry *procconf.ProcEntry) *Process {
	proc := &Process{
		cmd:        nil,
		procEntry:  procEntry,
		startTime:  time.Unix(0, 0),
		stopTime:   time.Unix(0, 0),
		state:      Stopped,
		inStart:    false,
		stopByUser: false,
		retryTimes: new(int32),
	}
	return proc
}

// Start 启动进程
func (that *Process) Start(wait bool) {
	loggerv2.Infof("try to start program %s", that.GetName())

	that.lock.Lock()
	if that.inStart {
		loggerv2.Infof("Don't start program again, program is already started,%s", that.GetName())
		that.lock.Unlock()
		return
	}
	that.inStart = true
	that.stopByUser = false
	that.lock.Unlock()

	var runCond *sync.Cond
	if wait {
		runCond = sync.NewCond(&sync.Mutex{})
		runCond.L.Lock()
	}

	go func() {
		for {
			that.run(func() {
				if wait {
					runCond.L.Lock()
					runCond.Signal()
					runCond.L.Unlock()
				}
			})

			// 如果启动进程是失败，一直重试，为了避免这种情况，暂停一会
			if time.Now().Unix()-that.startTime.Unix() < 2 {
				time.Sleep(5 * time.Second)
			}
			if that.stopByUser {
				loggerv2.Infof("%s Stopped by user, don't start it again", that.GetName())
				break
			}
			// 判断进程是否需要自动重启
			if !that.isAutoRestart() {
				loggerv2.Infof("Don't start the stopped program %s because its autorestart flag is false", that.GetName())
				break
			}
		}
		that.lock.Lock()
		that.inStart = false
		that.lock.Unlock()
	}()

	if wait {
		runCond.Wait()
		runCond.L.Unlock()
	}
}

// 启动进程
func (that *Process) run(finishCb func()) {
	that.lock.Lock()
	defer that.lock.Unlock()

	if that.isRunning() {
		loggerv2.Infof("Don't start program %s because it is running", that.GetName())
		finishCb()
		return
	}

	that.startTime = time.Now()
	atomic.StoreInt32(that.retryTimes, 0)
	//获取启动时间
	startSecs := that.procEntry.StartSecs()
	// 重启暂停时间
	restartPause := that.procEntry.RestartPause()

	var once sync.Once
	finishCbWrapper := func() {
		once.Do(finishCb)
	}

	// 进程没有超时，且没有被用户结束
	for !that.stopByUser {

		//如果进程启动失败，需要重试，则需要判断配置，重试启动是否需要间隔制定时间
		if restartPause > 0 && atomic.LoadInt32(that.retryTimes) != 0 {
			that.lock.Lock()
			loggerv2.Infof("don't restart the program %s, start it after ", that.GetName())
			time.Sleep(time.Duration(restartPause) * time.Second)
			that.lock.Unlock()
		}

		endTime := time.Now().Add(time.Duration(startSecs) * time.Second)
		//更新进程状态
		that.changeStateTo(Starting)
		// 启动次数+1
		atomic.AddInt32(that.retryTimes, 1)
		// 创建启动命令行
		err := that.createProgramCommand()
		if err != nil {
			that.failToStartProgram("fail to create program", finishCbWrapper)
			break
		}
		// 启动
		err = that.cmd.Start()
		if err != nil {
			// 重试次数已经大于设置中的最大重试次数
			if atomic.LoadInt32(that.retryTimes) >= int32(that.procEntry.StartRetries()) {
				that.failToStartProgram(fmt.Sprintf("fail to start program with error:%v", err), finishCbWrapper)
				break
			} else {
				loggerv2.Infof("program:%s fail to start program with error:%v", that.GetName(), err)
				that.changeStateTo(Backoff)
				continue
			}
		}
		if that.StdoutLog != nil {
			that.StdoutLog.SetPid(that.cmd.Process.Pid)
		}
		if that.StderrLog != nil {
			that.StderrLog.SetPid(that.cmd.Process.Pid)
		}
		monitorExited := int32(0)
		programExited := int32(0)
		if startSecs <= 0 {
			loggerv2.Infof("%s success to start program", that.GetName())
			that.changeStateTo(Running)
			go finishCbWrapper()
		} else {
			go func() {
				that.monitorProgramIsRunning(endTime, &monitorExited, &programExited)
				finishCbWrapper()
			}()
		}
		loggerv2.Debugf("%s wait program exit", that.GetName())
		that.lock.Unlock()
		that.waitForExit(int64(startSecs))
		atomic.StoreInt32(&programExited, 1)
		// wait for monitor thread exit
		for atomic.LoadInt32(&monitorExited) == 0 {
			time.Sleep(time.Duration(10) * time.Millisecond)
		}
		that.lock.Lock()

		// if the program still in running after startSecs
		if that.state == Running {
			that.changeStateTo(Exited)
			loggerv2.Infof("%s program exited", that.GetName())
			break
		} else {
			that.changeStateTo(Backoff)
		}
		if atomic.LoadInt32(that.retryTimes) >= int32(that.procEntry.StartRetries()) {
			that.failToStartProgram(fmt.Sprintf("fail to start program because retry times is greater than %d", that.procEntry.StartRetries), finishCbWrapper)
			break
		}
	}
}

// 创建程序的cmd对象
func (that *Process) createProgramCommand() (err error) {
	that.cmd, err = that.procEntry.CreateCommand()
	if err != nil {
		return err
	}
	if that.setUser() != nil {
		loggerv2.Errorf("fail to run as user %s", that.procEntry.User)
		return fmt.Errorf("fail to set user")
	}
	setDeathSig(that.cmd.SysProcAttr)
	that.setEnv()
	that.setDir()
	that.setLog()
	that.stdin, _ = that.cmd.StdinPipe()
	return nil
}

// Stop 停止进程
func (that *Process) Stop(wait bool) {
	that.lock.Lock()
	that.stopByUser = true
	isRunning := that.isRunning()
	that.lock.Unlock()
	if !isRunning {
		loggerv2.Infof("program %s is not running", that.GetName())
		return
	}
	loggerv2.Infof("stop the program %s", that.GetName())

	sigs := strings.Fields(that.procEntry.StopSignal())
	waitSecs := time.Duration(that.procEntry.StopWaitSecs(10)) * time.Second
	killWaitSecs := time.Duration(that.procEntry.KillWaitSecs(2)) * time.Second
	stopAsGroup := that.procEntry.StopAsGroup()
	killAsGroup := that.procEntry.KillAsGroup()
	if stopAsGroup && !killAsGroup {
		loggerv2.Error("Cannot set stopAsGroup=true and killAsGroup=false")
	}
	var stopped int32 = 0

	go func() {
		for i := 0; i < len(sigs) && atomic.LoadInt32(&stopped) == 0; i++ {
			// send signal to process
			sig, err := signals.ToSignal(sigs[i])
			if err != nil {
				continue
			}
			loggerv2.Infof("send stop signal %s to program %s", that.GetName(), sigs[i])
			_ = that.Signal(sig, stopAsGroup)
			endTime := time.Now().Add(waitSecs)
			//wait at most "stopwaitsecs" seconds for one signal
			for endTime.After(time.Now()) {
				//if it already exits
				if that.state != Starting && that.state != Running && that.state != Stopping {
					atomic.StoreInt32(&stopped, 1)
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
		if atomic.LoadInt32(&stopped) == 0 {
			loggerv2.Infof("force to kill the program%s", that.GetName())
			_ = that.Signal(syscall.SIGKILL, killAsGroup)
			killEndTime := time.Now().Add(killWaitSecs)
			for killEndTime.After(time.Now()) {
				//if it exits
				if that.state != Starting && that.state != Running && that.state != Stopping {
					atomic.StoreInt32(&stopped, 1)
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			atomic.StoreInt32(&stopped, 1)
		}
	}()
	if wait {
		for atomic.LoadInt32(&stopped) == 0 {
			time.Sleep(1 * time.Second)
		}
	}
}
