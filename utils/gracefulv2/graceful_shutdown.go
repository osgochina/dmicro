package gracefulv2

import (
	"context"
	"github.com/osgochina/dmicro/logger"
	"os"
	"time"
)

// Shutdown 执行进程关闭任务
func (that *Graceful) Shutdown(timeout ...time.Duration) {
	defer os.Exit(0)
	var isReboot = false
	if that.processStatus.Val() == statusActionRestarting {
		isReboot = true
		logger.Infof("平滑重启，正在结束父进程...")
	} else {
		logger.Infof("正在结束进程...")
	}
	that.processStatus.Set(statusActionShuttingDown)

	that.contextExec(timeout, "shutdown", func(ctxTimeout context.Context) <-chan struct{} {
		endCh := make(chan struct{})
		go func() {
			defer close(endCh)
			var graceful = true
			//当进程非重启状态时候，才需要执行清理动作
			if !isReboot {
				if err := that.firstSweep(); err != nil {
					logger.Errorf("[结束进程 - 执行前置方法失败] %s", err.Error())
					graceful = false
				}
			}
			graceful = that.shutdown(ctxTimeout, "shutdown") && graceful
			if graceful {
				logger.Info("进程结束了.")
			} else {
				logger.Info("进程结束了,但是非平滑模式.")
			}
		}()
		return endCh
	})
}

func (that *Graceful) ShutdownMasterWorker(timeout ...time.Duration) {
	defer os.Exit(0)
	logger.Infof("正在结束进程...")
	that.processStatus.Set(statusActionShuttingDown)

	that.contextExec(timeout, "shutdown", func(ctxTimeout context.Context) <-chan struct{} {
		endCh := make(chan struct{})
		go func() {
			defer close(endCh)
			var graceful = true
			if err := that.firstSweep(); err != nil {
				logger.Errorf("[结束进程 - 执行前置方法失败] %s", err.Error())
				graceful = false
			}
			graceful = that.shutdown(ctxTimeout, "shutdown") && graceful
			if graceful {
				logger.Info("进程结束了.")
			} else {
				logger.Info("进程结束了,但是非平滑模式.")
			}
		}()
		return endCh
	})
}

// 执行shutdown和reboot命令，并且计时，在规定的时候内为执行完收尾动作，则强制结束进程
func (that *Graceful) contextExec(timeout []time.Duration, action string, deferCallback func(ctxTimeout context.Context) <-chan struct{}) {
	if len(timeout) > 0 {
		that.SetShutdown(timeout[0], that.firstSweep, that.beforeExiting)
	}
	ctxTimeout, cancel := context.WithTimeout(context.Background(), that.shutdownTimeout)
	defer cancel()
	select {
	case <-ctxTimeout.Done():
		if err := ctxTimeout.Err(); err != nil {
			logger.Errorf("[%s-timeout] %s", action, err.Error())
		}
	case <-deferCallback(ctxTimeout):
	}
}

//执行后置函数
func (that *Graceful) shutdown(ctxTimeout context.Context, action string) bool {
	logger.Info("[结束进程中 - 正在执行后置函数]")
	if err := that.beforeExiting(); err != nil {
		logger.Errorf("[%s-beforeExiting] %v", action, err)
		return false
	}
	return true
}
