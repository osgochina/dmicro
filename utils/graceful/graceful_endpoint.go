package graceful

import (
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/netproto/quic"
	"github.com/osgochina/dmicro/logger"
	"net"
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
	if !that.isChild() {
		return
	}
	if that.model != GraceChangeProcess {
		return
	}
	pPid := syscall.Getppid()
	if pPid != 1 {
		if err := syscallKillSIGTERM(pPid); err != nil {
			logger.Errorf("子进程重启后向父进程发送信号失败，error: %s", err.Error())
			return
		}
		logger.Printf("平滑重启中,子进程[%d]已向父进程[%d]发送信号'SIGTERM'", syscall.Getpid(), pPid)
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

// 获取监听列表
func (that *graceful) getEndpointListenerFdMap() map[string]string {
	if that.model == GraceMasterWorker {
		return that.getEndpointListenerFdMasterWorker()
	}
	return that.getEndpointListenerFdMapChangeProcess()
}

// 获取父子进程模式下的监听fd列表
func (that *graceful) getEndpointListenerFdMapChangeProcess() map[string]string {
	if that.inheritedProcListener.Len() <= 0 {
		return nil
	}
	m := map[string]string{
		"tcp":  "",
		"quic": "",
	}
	that.inheritedProcListener.Iterator(func(_ int, v interface{}) bool {
		lis, ok := v.(net.Listener)
		if !ok {
			logger.Warningf("inheritedProcListener 不是 net.Listener类型")
			return true
		}
		quicLis, ok := v.(*quic.Listener)
		if ok {
			f, e := quicLis.PacketConn().(filer).File()
			if e != nil {
				logger.Error(e)
				return false
			}
			str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
			if len(m["quic"]) > 0 {
				m["quic"] += ","
			}
			m["quic"] += str
			return true
		}
		f, e := lis.(filer).File()
		if e != nil {
			logger.Error(e)
			return false
		}
		str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
		if len(m["tcp"]) > 0 {
			m["tcp"] += ","
		}
		m["tcp"] += str
		return true
	})
	return m
}

// 获取master worker 进程模型下的监听列表
func (that *graceful) getEndpointListenerFdMasterWorker() map[string]string {
	if that.inheritedProcListener.Len() <= 0 {
		return nil
	}
	m := map[string]string{
		"tcp":  "",
		"quic": "",
	}
	that.inheritedProcListener.Iterator(func(_ int, v interface{}) bool {
		lis, ok := v.(net.Listener)
		if !ok {
			logger.Warningf("inheritedProcListener 不是 net.Listener类型")
			return true
		}
		if that.mwListenAddr != nil {
			// 判断监听的是否是http协议。如果是http协议则不返回
			data := that.mwListenAddr.Get(lis.Addr().String())
			if d, ok := data.(InheritAddr); ok && (d.Network == "http" || d.Network == "https") {
				return true
			}
		}
		quicLis, ok := v.(*quic.Listener)
		if ok {
			f, e := quicLis.PacketConn().(filer).File()
			if e != nil {
				logger.Error(e)
				return false
			}
			str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
			if len(m["quic"]) > 0 {
				m["quic"] += ","
			}
			m["quic"] += str
			return true
		}
		f, e := lis.(filer).File()
		if e != nil {
			logger.Error(e)
			return true
		}
		str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
		if len(m["tcp"]) > 0 {
			m["tcp"] += ","
		}
		m["tcp"] += str
		return true
	})

	return m
}
