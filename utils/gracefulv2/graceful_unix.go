// +build !windows

package gracefulv2

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
func (that *Graceful) AddInherited(procFiles []*os.File, envs map[string]string) {
	for _, f := range procFiles {
		// 判断需要添加的文件句柄是否已经存在,不存在才能追加
		if that.inheritedProcFiles.Search(f) == -1 {
			that.inheritedProcFiles.Append(f)
		}
	}
	that.inheritedEnv.Sets(envs)
}

func (that *Graceful) GraceSignal() {
	// subscribe to SIGINT signals
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
		logger.Printf(`收到信号: %s`, sig.String())
		switch sig {
		// 强制关闭服务
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
			// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
			that.Shutdown(time.Second)
			return
		// 平滑的关闭服务
		case syscall.SIGTERM:
			// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
			that.Shutdown()
			return
		// 平滑重启服务
		case syscall.SIGUSR1:
			that.Reboot()
			return
		default:
		}
	}
}

// Reboot 开启优雅的重启流程
func (that *Graceful) Reboot(timeout ...time.Duration) {
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

//启动新的进程
func (that *Graceful) startProcess() (int, error) {
	var extraFiles []*os.File
	that.inheritedProcFiles.Iterator(func(k int, v interface{}) bool {
		extraFiles = append(extraFiles, v.(*os.File))
		return true
	})

	//获取进程启动的原始
	path := os.Args[0]
	err := genv.SetMap(that.inheritedEnv.Map())
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
	cmd.ExtraFiles = extraFiles
	cmd.Env = envs
	cmd.Dir = originalWD
	err = cmd.Start()
	if err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
