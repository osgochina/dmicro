package normal

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// Listen 监听
func Listen(nett, addr string) (net.Listener, error) {
	return globalInheritNet.Listen(nett, addr)
}

// ListenTCP 监听tcp协议
func ListenTCP(nett string, addr *net.TCPAddr) (*net.TCPListener, error) {
	return globalInheritNet.ListenTCP(nett, addr)
}

// ListenUnix 监听unix协议
func ListenUnix(nett string, addr *net.UnixAddr) (*net.UnixListener, error) {
	return globalInheritNet.ListenUnix(nett, addr)
}

// Append 追加监听句柄到活跃列表
func Append(ln net.Listener) error {
	return globalInheritNet.Append(ln)
}

// SetInherited 添加files列表到环境变量，让子进程继承，
//1. 只有在reboot使用
//2. 不支持windows系统
func SetInherited() error {
	return globalInheritNet.SetInherited()
}

// AddInheritedFunc 平滑重启的时候，会回调该方法，保存fd列表
func AddInheritedFunc(fn func([]net.Listener, map[string]string)) {
	globalInheritNet.addInherited = fn
}

// GetInheritedFunc 如果是平滑重启，可以获取到从父进程继承过来的fd列表
func GetInheritedFunc(fn func() []int) {
	globalInheritNet.getInherited = fn
}

var globalInheritNet = new(inheritNet)

type inheritNet struct {
	//继承过来的监听句柄列表
	inherited []net.Listener
	//当前进程活跃使用的监听句柄列表
	active      []net.Listener
	mutex       sync.Mutex
	inheritOnce sync.Once

	//传递需要继承的文件句柄列表方法
	addInherited func([]net.Listener, map[string]string)
	getInherited func() []int
}

// 从父进程继承的句柄初始化
//注意，这里使用了 sync.Once 逻辑，保证了仅能执行一次
func (that *inheritNet) inherit() error {
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
			l, e := net.FileListener(file)
			if e != nil {
				_ = file.Close()
				retErr = fmt.Errorf("error inheriting socket fd %d: %s", fd, e)
				return
			}
			if e = file.Close(); e != nil {
				retErr = fmt.Errorf("error closing inherited socket fd %d: %s", fd, e)
				return
			}
			that.inherited = append(that.inherited, l)
		}
	})
	return retErr
}

// Listen 监听
func (that *inheritNet) Listen(nett, addr string) (net.Listener, error) {
	switch nett {
	default:
		return nil, net.UnknownNetworkError(nett)
	case "tcp", "tcp4", "tcp6":
		tcpAddr, err := net.ResolveTCPAddr(nett, addr)
		if err != nil {
			return nil, err
		}
		return that.ListenTCP(nett, tcpAddr)
	case "unix", "unixpacket", "invalid_unix_net_for_test":
		unixAddr, err := net.ResolveUnixAddr(nett, addr)
		if err != nil {
			return nil, err
		}
		return that.ListenUnix(nett, unixAddr)
	}
}

// ListenTCP 监听tcp句柄
func (that *inheritNet) ListenTCP(nett string, addr *net.TCPAddr) (*net.TCPListener, error) {
	//初始化继承过来的句柄,只会初始化一次
	if err := that.inherit(); err != nil {
		return nil, err
	}
	that.mutex.Lock()
	defer that.mutex.Unlock()

	for i, l := range that.inherited {
		//如果继承过来的句柄变成了nil，则跳过
		if l == nil { // we nil used inherited listeners
			continue
		}
		//如果将要监听的地址已经在继承列表中，则直接返回该继承的句柄
		if isSameAddr(l.Addr(), addr) {
			that.inherited[i] = nil              //如果地址相同，则把改地址从继承列表拿出来使用
			that.active = append(that.active, l) //把继承列表中的地址拿出来放到已使用列表
			return l.(*net.TCPListener), nil
		}
	}
	//如果不在继承列表中，则直接新建一个监听,并把它加入到活跃列表
	l, err := net.ListenTCP(nett, addr)
	if err != nil {
		return nil, err
	}
	that.active = append(that.active, l)
	return l, nil
}

// ListenUnix 监听unix 文件类型的句柄
func (that *inheritNet) ListenUnix(nett string, addr *net.UnixAddr) (*net.UnixListener, error) {
	//初始化继承过来的句柄,只会初始化一次
	if err := that.inherit(); err != nil {
		return nil, err
	}
	that.mutex.Lock()
	defer that.mutex.Unlock()

	for i, l := range that.inherited {
		//如果继承过来的句柄变成了nil，则跳过
		if l == nil { // we nil used inherited listeners
			continue
		}
		//如果将要监听的地址已经在继承列表中，则直接返回该继承的句柄
		if isSameAddr(l.Addr(), addr) {
			that.inherited[i] = nil              //如果地址相同，则把改地址从继承列表拿出来使用
			that.active = append(that.active, l) //把继承列表中的地址拿出来放到已使用列表
			return l.(*net.UnixListener), nil
		}
	}
	//创建新鲜的 listener
	l, err := net.ListenUnix(nett, addr)
	if err != nil {
		return nil, err
	}
	that.active = append(that.active, l)
	return l, nil
}

// Append 追加监听句柄
func (that *inheritNet) Append(ln net.Listener) error {
	//初始化继承过来的句柄,只会初始化一次
	if err := that.inherit(); err != nil {
		return err
	}
	that.mutex.Lock()
	defer that.mutex.Unlock()

	//先从继承列表中查找
	for i, l := range that.inherited {
		if l == nil { // nil 表示已经使用了
			continue
		}
		//如果找到，则标记加入
		if isSameAddr(l.Addr(), ln.Addr()) {
			that.inherited[i] = nil
			that.active = append(that.active, l)
			return nil
		}
	}
	for _, l := range that.active {
		if l == nil {
			continue
		}
		//重复添加监听句柄
		if isSameAddr(l.Addr(), ln.Addr()) {
			return fmt.Errorf(" Re-register the listening port: network %s, address %s", ln.Addr().Network(), ln.Addr().String())
		}
	}
	that.active = append(that.active, ln)
	return nil
}

func (that *inheritNet) SetInherited() error {
	//获取所有正在使用的句柄
	listeners, err := that.activeListeners()
	if err != nil {
		return err
	}
	that.addInherited(listeners, nil)
	return nil
}

// 获取当前进程正在使用的监听句柄
func (that *inheritNet) activeListeners() ([]net.Listener, error) {
	that.mutex.Lock()
	defer that.mutex.Unlock()
	ls := make([]net.Listener, len(that.active))
	copy(ls, that.active)
	return ls, nil
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
