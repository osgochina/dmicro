package graceful

import (
	"github.com/gogf/gf/container/gset"
	"github.com/osgochina/dmicro/logger"
	"syscall"
)

// AddEndpoint 启动了endpoint后，添加到列表，在重启的时候方便轮训close
func (that *graceful) AddEndpoint(e interface{}) {
	that.endpointList.Add(e)
}

// DeleteEndpoint 当endpoint关闭后，从列表中删除它
func (that *graceful) DeleteEndpoint(e interface{}) {
	that.endpointList.Remove(e)
}

// 父子进程模型下，当子进程启动成功，发送信号通知父进程
func (that *graceful) onStart() {
	//非子进程，则什么都不走
	if that.IsChild() == false {
		return
	}
	if that.model != GraceChangeProcess {
		return
	}
	pPid := syscall.Getppid()
	if pPid != 1 {
		if err := syscallKillSIGTERM(pPid); err != nil {
			logger.Errorf("[reboot-killOldProcess] %s", err.Error())
			return
		}
		logger.Infof("平滑重启中,子进程[%d]已向父进程[%d]发送信号'SIGTERM'", syscall.Getpid(), pPid)
	}
}

func (that *graceful) SetShutdownEndpoint(callback func(*gset.Set) error) {
	that.shutdownCallback = callback
}

// 阻塞的等待Endpoint结束
func (that *graceful) shutdownEndpoint() error {
	if that.shutdownCallback == nil {
		return nil
	}
	return that.shutdownCallback(that.endpointList)
}
