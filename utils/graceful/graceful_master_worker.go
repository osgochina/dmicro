package graceful

import (
	"crypto/tls"
	"net"
	"os"
	"sync"
	"time"
)

type ListenAddr struct {
	Addr      net.Addr
	TlsConfig *tls.Config
}

type MasterWorkerGraceful struct {
	Graceful
	// 监听的信号
	signal chan os.Signal
	// 平滑重启的端口列表，支持:"127.0.0.1:8080",":8081","/var/local.sock"等
	listenAddr []*ListenAddr
	// 重启的等候时间，超过该时间未完成，则强制发送kill 9
	shutdownTimeout time.Duration
	// 锁
	locker        sync.Mutex
	firstSweep    func() error
	beforeExiting func() error
}

func NewMasterWorkerGraceful() *MasterWorkerGraceful {
	graceful := &MasterWorkerGraceful{
		signal:          make(chan os.Signal),
		shutdownTimeout: MinShutdownTimeout,
		firstSweep: func() error {
			return nil
		},
		beforeExiting: func() error {
			return nil
		},
	}
	graceful.parentAddrList = make(map[string]map[string][]string, 2)
	return graceful
}

func (that *MasterWorkerGraceful) SetListenAddr(listenAddr []*ListenAddr) {
	that.listenAddr = listenAddr
}

func (that *MasterWorkerGraceful) Listens() {
	//for _,addr := range that.listenAddr{
	//	//inherit.Listen(addr.Addr,addr.TlsConfig)
	//}

}
