// +build !windows

package graceful

import (
	"context"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/osgochina/dmicro/logger"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

//进程启动时候的原始路径
var originalWD, _ = os.Getwd()

func (that *Graceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	sig := <-that.signal
	signal.Stop(that.signal)
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		that.Shutdown()
	case syscall.SIGUSR2:
		that.Reboot()
	}
}

// Reboot 开启优雅的重启流程
func (that *Graceful) Reboot(timeout ...time.Duration) {
	defer os.Exit(0)
	logger.Info("rebooting process...")
	var (
		//ppid     = os.Getppid()
		graceful = true
	)

	that.contextExec(timeout, "reboot", func(ctxTimeout context.Context) <-chan struct{} {
		endCh := make(chan struct{})

		go func() {
			defer close(endCh)
			var reboot = true
			if err := that.firstSweep(); err != nil {
				logger.Infof("[reboot-firstSweep] %s", err.Error())
				graceful = false
			}

			//启动新的进程
			_, err := that.startProcess()
			if err != nil {
				logger.Infof("[reboot-startNewProcess] %s", err.Error())
				reboot = false
			}

			// 关闭当前进程
			graceful = that.shutdown(ctxTimeout, "reboot") && graceful
			if !reboot {
				if graceful {
					logger.Warning("process reboot failed, but shut down gracefully!")
				} else {
					logger.Warning("process reboot failed, and did not shut down gracefully!")
				}
				os.Exit(-1)
			}
		}()

		return endCh
	})

	////如果父进程不是初始进程1，则关闭父进程,如果通过supervisor启动，会有问题，所以去掉该逻辑
	//if ppid != 1 {
	//	if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
	//		that.logger.Errorf("[reboot-killOldProcess] %s", err.Error())
	//		graceful = false
	//	}
	//}
	if graceful {
		logger.Infof("process is rebooted gracefully.")
	} else {
		logger.Infof("process is rebooted, but not gracefully.")
	}
}

// AddInherited 添加需要给重启后新进程继承的文件句柄和环境变量
func (that *Graceful) AddInherited(procFiles []*os.File, envs map[string]string) {
	that.locker.Lock()
	defer that.locker.Unlock()
	for _, f := range procFiles {
		var had bool
		for _, ff := range that.inheritedProcFiles {
			if ff == f {
				had = true
				break
			}
		}
		if !had {
			that.inheritedProcFiles = append(that.inheritedProcFiles, f)
		}
	}
	for k, v := range envs {
		that.inheritedEnv[k] = v
	}
}

//启动新的进程
func (that *Graceful) startProcess() (int, error) {

	//关闭当前进程的指定默认句柄，为了给新进程让路
	for i, f := range that.inheritedProcFiles {
		if i >= that.defaultInheritedProcFilesLen {
			defer func() {
				_ = f.Close()
			}()
		}
	}

	//获取进程启动的原始
	path := gfile.SelfPath()
	err := genv.SetMap(that.inheritedEnv)
	if err != nil {
		return 0, err
	}
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	envs := genv.All()
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = that.inheritedProcFiles
	cmd.Env = envs
	cmd.Dir = originalWD
	err = cmd.Start()
	if err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
