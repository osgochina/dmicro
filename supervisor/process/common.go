package process

import (
	"fmt"
	loggerv2 "github.com/osgochina/dmicro/logger"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

// 判断进程是否在运行
func (that *Process) isRunning() bool {
	if that.cmd != nil && that.cmd.Process != nil {
		if runtime.GOOS == "windows" {
			proc, err := os.FindProcess(that.cmd.Process.Pid)
			return proc != nil && err == nil
		}
		return that.cmd.Process.Signal(syscall.Signal(0)) == nil
	}
	return false
}

// 在supervisord启动的时候也自动启动
func (that *Process) isAutoStart() bool {
	return that.procEntry.AutoStart()
}

// 设置进程的运行用户
func (that *Process) setUser() error {
	userName := that.procEntry.User()
	if len(userName) == 0 {
		return nil
	}

	//check if group is provided
	pos := strings.Index(userName, ":")
	groupName := ""
	if pos != -1 {
		groupName = userName[pos+1:]
		userName = userName[0:pos]
	}
	u, err := user.Lookup(userName)
	if err != nil {
		return err
	}
	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return err
	}
	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil && groupName == "" {
		return err
	}
	if groupName != "" {
		g, err := user.LookupGroup(groupName)
		if err != nil {
			return err
		}
		gid, err = strconv.ParseUint(g.Gid, 10, 32)
		if err != nil {
			return err
		}
	}
	that.cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid), NoSetGroups: true}
	return nil
}

// 设置进程运行的环境变量
func (that *Process) setEnv() {

	if len(that.procEntry.Environment()) != 0 {
		that.cmd.Env = os.Environ()
		for k, v := range that.procEntry.Environment() {
			that.cmd.Env = append(that.cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	} else {
		that.cmd.Env = os.Environ()
	}
}

// 设置进程的运行目录
func (that *Process) setDir() {
	dir := that.procEntry.Directory()
	if dir != "" {
		that.cmd.Dir = dir
	}
}

func (that *Process) setLog() {
	that.StdoutLog = that.createStdoutLogger()
	that.cmd.Stdout = that.StdoutLog
	if that.procEntry.RedirectStderr() {
		that.StderrLog = that.StdoutLog
	} else {
		that.StderrLog = that.createStderrLogger()
	}
	that.cmd.Stderr = that.StderrLog
}

func (that *Process) failToStartProgram(reason string, finishCb func()) {
	loggerv2.Errorf("%s program:%s", reason, that.GetName())
	that.changeStateTo(Fatal)
	finishCb()
}

func (that *Process) monitorProgramIsRunning(endTime time.Time, monitorExited *int32, programExited *int32) {
	// if time is not expired
	for time.Now().Before(endTime) && atomic.LoadInt32(programExited) == 0 {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	atomic.StoreInt32(monitorExited, 1)

	that.lock.Lock()
	defer that.lock.Unlock()
	// if the program does not exit
	if atomic.LoadInt32(programExited) == 0 && that.state == Starting {
		loggerv2.Infof("success to start program %s", that.GetName())
		that.changeStateTo(Running)
	}
}

// 判断进程是否需要自动重启
func (that *Process) isAutoRestart() bool {
	autoRestart := that.procEntry.AutoReStart("unexpected")

	if autoRestart == "false" {
		return false
	} else if autoRestart == "true" {
		return true
	} else {
		that.lock.RLock()
		defer that.lock.RUnlock()
		if that.cmd != nil && that.cmd.ProcessState != nil {
			exitCode, err := that.getExitCode()
			//If unexpected, the process will be restarted when the program exits
			//with an exit code that is not one of the exit codes associated with
			//this process’ configuration (see exitcodes).
			return err == nil && !that.inExitCodes(exitCode)
		}
	}
	return false
}

func (that *Process) waitForExit(startSecs int64) {
	_ = that.cmd.Wait()
	if that.cmd.ProcessState != nil {
		loggerv2.Infof(" program %s stopped with status:%v", that.GetName(), that.cmd.ProcessState)
	} else {
		loggerv2.Infof("program %s stopped", that.GetName())
	}
	that.lock.Lock()
	defer that.lock.Unlock()
	that.stopTime = time.Now()
	if that.StdoutLog != nil {
		_ = that.StdoutLog.Close()
	}
	if that.StderrLog != nil {
		_ = that.StderrLog.Close()
	}
}
