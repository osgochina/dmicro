package gracefulv2

import (
	"github.com/gogf/gf/container/gset"
	"github.com/osgochina/dmicro/logger"
	"net"
	"syscall"
)

func (that *Graceful) AddEndpoint(e interface{}) {
	that.endpointList.Add(e)
}

func (that *Graceful) DeleteEndpoint(e interface{}) {
	that.endpointList.Remove(e)
}

func (that *Graceful) OnListen(addr net.Addr) {
	//非子进程，则什么都不走
	if that.IsChild() == false {
		return
	}
	if that.model != GracefulChangeProcess {
		return
	}
	pPid := syscall.Getppid()
	if pPid != 1 {
		if err := SyscallKillSIGTERM(pPid); err != nil {
			logger.Errorf("[reboot-killOldProcess] %s", err.Error())
			return
		}
		logger.Infof("平滑重启中,子进程[%d]已向父进程[%d]发送信号'SIGTERM'", syscall.Getpid(), pPid)
	}
}

func (that *Graceful) SetShutdownEndpoint(callback func(*gset.Set) error) {
	that.shutdownCallback = callback
}

// 阻塞的等待Endpoint结束
func (that *Graceful) shutdownEndpoint() error {
	if that.shutdownCallback == nil {
		return nil
	}
	return that.shutdownCallback(that.endpointList)
}
