package graceful

import (
	"github.com/osgochina/dmicro/logger"
	"os/signal"
	"syscall"
)

func (that *MasterWorkerGraceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-that.signal
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			signal.Reset(syscall.SIGINT, syscall.SIGTERM)
			logger.Infof("收到了关闭信号%v", sig)
			//判断当前进程是Master进程还是Worker进程，如果是Master进程，则发送信号给worker进程，当监听到worker进程退出后，再自己退出。
			// 如果是worker进程，则处理完任务后,自己退出。
		case syscall.SIGUSR1:
			signal.Reset(syscall.SIGUSR1)
			logger.Infof("收到了重启信号%v", sig)
			//判断当前进程是Master进程还是Worker进程，
			//如果是Master进程，先启动一个新的worker进程，启动成功后该worker进程会发送信号SIGUSR2，
			//master进程收到信号后再则发送信号SIGTERM给老的worker进程，
			//当监听到老的worker进程退出，表示重启成功
			//如果是worker进程,则忽略该信号
		case syscall.SIGUSR2:
			signal.Reset(syscall.SIGUSR2)
			logger.Infof("收到了重启信号%v", sig)
			//判断当前进程是Master进程还是Worker进程，
			//如果是Master进程，则表示该worker进程已经启动成功
			//如果是worker进程,则忽略该信号
		}
	}
}
