package kcp

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

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

var globalInheritKCP = new(inheritKCP)

type inheritKCP struct {
	inherited   []*net.UDPConn
	active      []*Listener
	mutex       sync.Mutex
	inheritOnce sync.Once

	//传递需要继承的文件句柄列表方法
	addInherited func([]*Listener, map[string]string)
	getInherited func() []int
}

func (that *inheritKCP) inheritedListen(network string, udpAddr *net.UDPAddr, tlsConf *tls.Config, dataShards, parityShards int) (*Listener, error) {
	if err := that.inherit(); err != nil {
		return nil, err
	}

	that.mutex.Lock()
	defer that.mutex.Unlock()

	var udpConn *net.UDPConn

	// look for an inherited listener
	for i, conn := range that.inherited {
		if conn == nil { // we nil used inherited listeners
			continue
		}
		if isSameAddr(conn.LocalAddr(), udpAddr) {
			that.inherited[i] = nil
			udpConn = conn
		}
	}

	// make a fresh listener
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

func (that *inheritKCP) setInherited() error {
	listeners, err := that.activeListeners()
	if err != nil {
		return err
	}
	that.addInherited(listeners, nil)

	return nil
}

func (that *inheritKCP) activeListeners() ([]*Listener, error) {
	that.mutex.Lock()
	defer that.mutex.Unlock()
	ls := make([]*Listener, len(that.active))
	copy(ls, that.active)
	return ls, nil
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

//判断两个地址是否相同
func isSameAddr(addOne, addTwo net.Addr) bool {
	if addOne.Network() != addTwo.Network() {
		return false
	}
	addOneStr := addOne.String()
	addTwoStr := addTwo.String()

	if addOneStr == addTwoStr {
		return true
	}
	//去掉地址上的ipv6前缀
	const ipv6prefix = "[::]"
	addOneStr = strings.TrimPrefix(addOneStr, ipv6prefix)
	addTwoStr = strings.TrimPrefix(addTwoStr, ipv6prefix)

	//去掉地址上的ipv4前缀
	const ipv4prefix = "0.0.0.0"
	addOneStr = strings.TrimPrefix(addOneStr, ipv4prefix)
	addTwoStr = strings.TrimPrefix(addTwoStr, ipv4prefix)

	//判断去掉前缀后的地址是否相等
	return addOneStr == addTwoStr
}
