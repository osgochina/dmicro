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

var originalWD, _ = os.Getwd()

// AddInherited 添加需要给重启后新进程继承的文件句柄和环境变量
func (that *graceful) AddInherited(procFiles []*os.File, envs map[string]string) {
	for _, f := range procFiles {
		// 判断需要添加的文件句柄是否已经存在,不存在才能追加
		if that.inheritedProcFiles.Search(f) == -1 {
			that.inheritedProcFiles.Append(f)
		}
	}
	that.inheritedEnv.Sets(envs)
}

// GraceSignal 监听信号
func (that *graceful) GraceSignal() {
	if that.model == GraceChangeProcess {
		that.graceSignalGracefulChangeProcess()
		return
	}
	if that.model == GraceMasterWorker {
		that.graceSignalGracefulMW()
		return
	}
}

// 父子进程模式平滑重启
func (that *graceful) graceSignalGracefulChangeProcess() {
	signal.Notify(
		that.signal,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	for {
		sig := <-that.signal
		logger.Infof(`收到信号: %s`, sig.String())
		switch sig {
		// 强制关闭服务
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
			signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
			// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
			that.Shutdown(time.Second)
			continue
		// 平滑的关闭服务
		case syscall.SIGTERM:
			signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
			// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
			that.Shutdown()
			continue
		// 平滑重启服务
		case syscall.SIGUSR1:
			signal.Reset(syscall.SIGUSR1)
			that.Reboot()
			continue
		default:
		}
	}
}

// MasterWorker模式平滑重启
func (that *graceful) graceSignalGracefulMW() {
	if that.isChild() {
		signal.Notify(
			that.signal,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGABRT,
		)
		for {
			sig := <-that.signal
			logger.Infof(`收到信号: %s`, sig.String())
			switch sig {
			// 强制关闭服务
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
				that.Shutdown(time.Second)
				continue
			// 平滑的关闭服务
			case syscall.SIGTERM:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
				that.Shutdown()
				continue
			default:
			}
		}
	} else {
		signal.Notify(
			that.signal,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGABRT,
			syscall.SIGUSR1,
			syscall.SIGUSR2,
		)
		for {
			sig := <-that.signal
			logger.Infof(`收到信号: %s`, sig.String())
			switch sig {
			// 强制关闭服务
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
				that.shutdownMaster()
				continue
			// 平滑的关闭服务
			case syscall.SIGTERM:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
				that.shutdownMaster()
				continue
			// 平滑重启服务
			case syscall.SIGUSR1:
				signal.Reset(syscall.SIGUSR1)
				that.Reboot()
				continue
			default:
			}
		}
	}
}

// Reboot 开启优雅的重启流程
func (that *graceful) Reboot(timeout ...time.Duration) {
	if that.model == GraceMasterWorker {
		that.rebootMasterWorker()
		return
	}
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

// master worker模式重启，就是对子进程发送退出信号
func (that *graceful) rebootMasterWorker() {
	pid := that.mwChildCmd.Process.Pid
	logger.Infof(`向子进程: %d 发送信号SIGTERM`, pid)
	_ = syscallKillSIGTERM(pid)
}

//master worker进程模式，退出master进程方法
func (that *graceful) shutdownMaster() {
	defer os.Exit(0)
	pid := that.mwChildCmd.Process.Pid
	logger.Infof(`向子进程: %d 发送信号SIGTERM`, pid)
	_ = syscallKillSIGTERM(pid)
}

//启动新的进程
func (that *graceful) startProcess() (*exec.Cmd, error) {
	var extraFiles []*os.File
	that.inheritedProcFiles.Iterator(func(k int, v interface{}) bool {
		extraFiles = append(extraFiles, v.(*os.File))
		return true
	})

	//获取进程启动的原始
	path := os.Args[0]
	err := genv.SetMap(that.inheritedEnv.Map())
	if err != nil {
		return nil, err
	}

	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	envs := genv.All()
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = extraFiles
	cmd.Env = envs
	cmd.Dir = originalWD
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// 发送结束信号给进程
func syscallKillSIGTERM(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
