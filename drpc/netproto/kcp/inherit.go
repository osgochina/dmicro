package kcp

import (
	"crypto/tls"
	"fmt"
	"github.com/osgochina/dmicro/utils"
	"net"
	"os"
	"sync"
)

// InheritedListen 使用kcp协议启动监听，需要先判断是否是继承过来的端口
// network: 网络类型，在这里只能是udp协议，但是区分udp，udp4,udp6
// laddr: 本地要监听的数据包地址端口
// tlsConf: 证书信息
// dataShards,parityShards: FEC前向纠错算法的数据冗余度配置，这两个参数决定了fec的冗余度，冗余度越大抗丢包性就越强。
func InheritedListen(network, laddr string, tlsConf *tls.Config, dataShards, parityShards int) (net.Listener, error) {
	udpAddr, err := net.ResolveUDPAddr(network, laddr)
	if err != nil {
		return nil, err
	}
	return globalInheritKCP.inheritedListen(network, udpAddr, tlsConf, dataShards, parityShards)
}

// SetInherited 添加files列表到环境变量，让子进程继承，
//1. 只有在reboot使用
//2. 不支持windows系统
func SetInherited() error {
	return globalInheritKCP.setInherited()
}

// AddInheritedFunc 平滑重启的时候，会回调该方法，保存fd列表
func AddInheritedFunc(fn func([]*Listener, map[string]string)) {
	globalInheritKCP.addInherited = fn
}

// GetInheritedFunc 如果是平滑重启，可以获取到从父进程继承过来的fd列表
func GetInheritedFunc(fn func() []int) {
	globalInheritKCP.getInherited = fn
}

//#########################################以上的接口暴露给外部，以下接口内部使用#############################################

var globalInheritKCP = new(inheritKCP)

type inheritKCP struct {
	inherited   []*net.UDPConn
	active      []*Listener
	mutex       sync.Mutex
	inheritOnce sync.Once

	//传递需要继承的文件句柄列表方法
	addInherited func([]*Listener, map[string]string)
	// 获取从父进程继承过来的句柄列表
	getInherited func() []int
}

// 添加files列表到环境变量，让子进程继承，
//1. 只有在reboot使用
//2. 不支持windows系统
func (that *inheritKCP) setInherited() error {
	listeners, err := that.activeListeners()
	if err != nil {
		return err
	}
	if that.addInherited != nil {
		that.addInherited(listeners, nil)
	}
	return nil
}

// 获取当前进程正在使用的监听句柄
func (that *inheritKCP) activeListeners() ([]*Listener, error) {
	that.mutex.Lock()
	defer that.mutex.Unlock()
	ls := make([]*Listener, len(that.active))
	copy(ls, that.active)
	return ls, nil
}

// InheritedListen 监听地址，需要先判断是否有继承过来的句柄，
// 如果是子进程，并且已经继承了该地址的监听，则返回已监听的句柄
// 如果未发现该地址的监听，则重新创建监听
func (that *inheritKCP) inheritedListen(network string, udpAddr *net.UDPAddr, tlsConf *tls.Config, dataShards, parityShards int) (*Listener, error) {
	//初始化继承过来的句柄,只会初始化一次
	if err := that.inherit(); err != nil {
		return nil, err
	}

	that.mutex.Lock()
	defer that.mutex.Unlock()

	var udpConn *net.UDPConn

	// look for an inherited listener
	for i, conn := range that.inherited {
		//如果继承过来的句柄变成了nil，则跳过
		if conn == nil {
			continue
		}
		//如果将要监听的地址已经在继承列表中，则直接返回该继承的句柄
		if utils.IsSameAddr(conn.LocalAddr(), udpAddr) {
			that.inherited[i] = nil //如果地址相同，则把改地址从继承列表拿出来使用
			udpConn = conn
		}
	}

	//如果不在继承列表中，则直接新建一个监听,并把它加入到活跃列表
	var l *Listener
	var err error
	if udpConn == nil {
		l, err = ListenUDPAddr(network, udpAddr, tlsConf, dataShards, parityShards)
	} else {
		l, err = Listen(udpConn, tlsConf, dataShards, parityShards)
	}
	if err != nil {
		return nil, err
	}
	that.active = append(that.active, l)
	return l, nil
}

// 从父进程继承的句柄初始化
//注意，这里使用了 sync.Once 逻辑，保证了仅能执行一次
func (that *inheritKCP) inherit() error {
	var retErr error
	that.inheritOnce.Do(func() {
		that.mutex.Lock()
		defer that.mutex.Unlock()
		if that.getInherited == nil {
			return
		}
		fds := that.getInherited()
		if len(fds) <= 0 {
			return
		}

		for _, fd := range fds {
			file := os.NewFile(uintptr(fd), "listener")
			conn, err := net.FilePacketConn(file)
			if err != nil {
				file.Close()
				retErr = fmt.Errorf("error inheriting socket fd %d: %s", fd, err)
				return
			}
			if err := file.Close(); err != nil {
				retErr = fmt.Errorf("error closing inherited socket fd %d: %s", fd, err)
				return
			}
			if udpConn, ok := conn.(*net.UDPConn); ok {
				that.inherited = append(that.inherited, udpConn)
			}
		}
	})
	return retErr
}
