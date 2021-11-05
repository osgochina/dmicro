// +build !windows

package graceful

import (
	"context"
	"github.com/gogf/gf/os/genv"
	"github.com/osgochina/dmicro/logger"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

//进程启动时候的原始路径
var originalWD, _ = os.Getwd()
var isReboot = false

func (that *ChangeProcessGraceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	for {
		sig := <-that.signal
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			signal.Reset(syscall.SIGINT, syscall.SIGTERM)
			logger.Infof("收到了关闭信号%v", sig)
			that.Shutdown()
		case syscall.SIGUSR1:
			signal.Reset(syscall.SIGUSR1)
			logger.Infof("收到了重启信号%v", sig)
			isReboot = true
			that.Reboot()
		}
	}
}

// Reboot 开启优雅的重启流程
func (that *ChangeProcessGraceful) Reboot(timeout ...time.Duration) {
	//defer os.Exit(0)
	logger.Info("平滑重启中...")

	that.contextExec(timeout, "reboot", func(ctxTimeout context.Context) <-chan struct{} {
		endCh := make(chan struct{})

		go func() {
			defer close(endCh)
			if err := that.firstSweep(); err != nil {
				logger.Infof("[平滑重启中 - 执行前置方法失败] %s", err.Error())
				os.Exit(-1)
			}

			//启动新的进程
			_, err := that.startProcess()
			// 启动新的进程失败，则表示该进程有问题，直接错误退出
			if err != nil {
				logger.Infof("[平滑重启中 - 启动新的进程失败] %s", err.Error())
				os.Exit(-1)
			}
		}()
		return endCh
	})
	logger.Infof("进程已进行平滑重启,等待子进程的信号...")
}

//// AddInherited 添加需要给重启后新进程继承的文件句柄和环境变量
//func (that *ChangeProcessGraceful) AddInherited(procFiles []*os.File, envs map[string]string) {
//	that.locker.Lock()
//	defer that.locker.Unlock()
//	for _, f := range procFiles {
//		var had bool
//		for _, ff := range that.inheritedProcFiles {
//			if ff == f {
//				had = true
//				break
//			}
//		}
//		if !had {
//			that.inheritedProcFiles = append(that.inheritedProcFiles, f)
//		}
//	}
//	for k, v := range envs {
//		that.inheritedEnv[k] = v
//	}
//}

//启动新的进程
func (that *ChangeProcessGraceful) startProcess() (int, error) {

	//关闭当前进程的指定默认句柄，为了给新进程让路
	for i, f := range that.inheritedProcFiles {
		if i >= that.defaultInheritedProcFilesLen {
			defer func() {
				_ = f.Close()
			}()
		}
	}
	//获取进程启动的原始
	path := os.Args[0]
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

// SyscallKillSIGTERM 发送结束信号给进程
func SyscallKillSIGTERM(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
