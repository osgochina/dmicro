package dserver

import (
	"context"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/supervisor/process"
	"os"
	"time"
)

// 单进程模式退出
func (that *graceful) shutdownSingle(timeout ...time.Duration) {
	pid := os.Getpid()
	defer os.Exit(0)
	var isReboot = false
	if that.processStatus.Val() == statusActionRestarting {
		isReboot = true
		logger.Printf("进程:%d,平滑重启，正在结束父进程...", pid)
	} else {
		logger.Printf("进程:%d,正在结束...", pid)
	}
	that.processStatus.Set(statusActionShuttingDown)

	that.contextExec(timeout, "shutdown", func(ctxTimeout context.Context) <-chan struct{} {
		endCh := make(chan struct{})
		go func() {
			defer close(endCh)
			var g = true
			//当进程非重启状态时候，才需要执行清理动作
			if !isReboot {
				if err := that.firstSweep(); err != nil {
					logger.Errorf("进程:%d 结束中 - 执行前置方法失败，error: %s", pid, err.Error())
					g = false
				}
			}
			g = that.callBeforeExiting(ctxTimeout, "shutdown") && g
			if g {
				logger.Printf("进程:%d 结束了.", pid)
			} else {
				logger.Printf("进程:%d 结束了,但是非平滑模式.", pid)
			}
		}()
		return endCh
	})
}

// 执行shutdown和reboot命令，并且计时，在规定的时候内为执行完收尾动作，则强制结束进程
func (that *graceful) contextExec(timeout []time.Duration, action string, deferCallback func(ctxTimeout context.Context) <-chan struct{}) {
	if len(timeout) > 0 {
		that.setShutdown(timeout[0], that.firstSweep, that.beforeExiting)
	}
	pid := os.Getpid()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), that.shutdownTimeout)
	defer cancel()
	select {
	case <-ctxTimeout.Done():
		if err := ctxTimeout.Err(); err != nil {
			logger.Errorf("进程:%d,处理 %s 超时 %s", pid, action, err.Error())
		}
	case <-deferCallback(ctxTimeout):
	}
}

//执行后置函数
func (that *graceful) callBeforeExiting(ctxTimeout context.Context, action string) bool {
	pid := os.Getpid()
	logger.Printf("进程:%d 正在结束中 - 正在执行后置函数", pid)
	// 这里实现的有问题，并不能控制超时，后期完善
	select {
	case <-ctxTimeout.Done():
		return false
	default:
		if err := that.beforeExiting(); err != nil {
			logger.Errorf("进程:%d [%s-后置函数执行失败] error:%v", pid, action, err)
			return false
		}
	}
	return true
}

// 多进程模式子进程退出
func (that *graceful) shutdownMultiChild(timeout ...time.Duration) {
	pid := os.Getpid()
	defer os.Exit(0)
	logger.Printf("进程:%d,正在退出...", pid)
	that.processStatus.Set(statusActionShuttingDown)

	that.contextExec(timeout, "shutdown", func(ctxTimeout context.Context) <-chan struct{} {
		endCh := make(chan struct{})
		go func() {
			defer close(endCh)
			var g = true
			if err := that.firstSweep(); err != nil {
				logger.Errorf("进程:%d 结束中 - 执行前置方法失败，error: %s", pid, err.Error())
				g = false
			}
			g = that.callBeforeExiting(ctxTimeout, "shutdown") && g
			if g {
				logger.Printf("进程:%d 结束了.", pid)
			} else {
				logger.Printf("进程:%d 结束了,但是非平滑模式.", pid)
			}
		}()
		return endCh
	})
}

//master worker进程模式，退出master进程方法
func (that *graceful) shutdownMultiMaster() {
	defer os.Exit(0)
	that.dServer.manager.ForEachProcess(func(p *process.Process) {
		p.Stop(true)
	})
}
