package dserver

import (
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/netproto/kcp"
	"github.com/osgochina/dmicro/drpc/netproto/normal"
	"github.com/osgochina/dmicro/drpc/netproto/quic"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils"
	"github.com/osgochina/dmicro/utils/errors"
	"github.com/osgochina/dmicro/utils/signals"
	"net"
	"os"
	"syscall"
	"time"
)

// InheritAddr master进程需要监听的配置
type InheritAddr struct {
	Network   string
	Host      string
	Port      int
	TlsConfig *tls.Config
	// ghttp服务专用
	ServerName string
}

// 优雅重启
type graceful struct {
	// server对象
	dServer *DServer
	// 是否是子进程
	childBool bool
	// 当前进程的状态
	processStatus *gtype.Int

	//监听的信号
	signal chan os.Signal

	// 将要被子进程继承的环境变量
	inheritedEnv *gmap.StrStrMap

	// 将要被子进程继承的监听列表
	inheritedProcListener *garray.Array

	// 进程收到退出或重启信号后，需要执行的方法
	firstSweep func() error
	// 进程真正退出前需要执行的方法
	beforeExiting func() error
	// 退出与重启流程最大等待时间
	shutdownTimeout time.Duration

	// 监听的句柄需要被子进程继承，需要设置提取这些句柄，并把它们设置为环境变量
	// 该方法留给业务设置，在结束进程的时候调用
	setInherited func() error

	//  多进程模式下模式 监听的地址列表
	listenAddrMap *gmap.StrAnyMap
}

//当前进程的状态
const (
	// 初始状态
	statusActionNone = 0
	// 进程在重启中
	statusActionRestarting = 1
	// 进程正在结束中
	statusActionShuttingDown = 2
)

const (
	// 当前是否是在子进程
	isChildKey = "GRACEFUL_IS_CHILD"
	// 父进程的监听列表
	parentAddrKey = "GRACEFUL_INHERIT_LISTEN_PARENT_ADDR"
	// gf框架的ghttp服务平滑重启key
	adminActionReloadEnvKey = "GF_SERVER_RELOAD"
)

//  进程收到结束或重启信号后，存活的最大时间
const minShutdownTimeout = 15 * time.Second

// 新建graceful对象
func newGraceful(server *DServer) *graceful {
	return &graceful{
		dServer:               server,
		processStatus:         gtype.NewInt(statusActionNone),
		signal:                make(chan os.Signal),
		inheritedEnv:          gmap.NewStrStrMap(true),
		inheritedProcListener: garray.NewArray(true),
		firstSweep:            func() error { return nil },
		beforeExiting:         func() error { return nil },
		childBool:             genv.GetVar(isChildKey, false).Bool(),
	}
}

// 父子进程模型下，当子进程启动成功，发送信号通知父进程
func (that *graceful) onStart() {
	//非子进程，则什么都不做
	if !that.isChild() {
		return
	}
	pPid := syscall.Getppid()
	if pPid != 1 {
		if err := signals.KillPid(pPid, signals.ToSignal("SIGTERM"), false); err != nil {
			logger.Errorf("子进程重启后向父进程发送信号失败，error: %s", err.Error())
			return
		}
		logger.Printf("平滑重启中,子进程[%d]已向父进程[%d]发送信号'SIGTERM'", syscall.Getpid(), pPid)
	}
}

// SetShutdown 设置退出的基本参数
// timeout 传入负数，表示永远不过期
// firstSweepFunc 收到退出或重启信号后，立刻需要执行的方法
// beforeExitingFunc 处理好收尾动作后，真正要退出了执行的方法
func (that *graceful) setShutdown(timeout time.Duration, firstSweepFunc, beforeExitingFunc func() error) {
	if timeout < 0 {
		that.shutdownTimeout = 1<<63 - 1
	} else if timeout < minShutdownTimeout {
		that.shutdownTimeout = minShutdownTimeout
	} else {
		that.shutdownTimeout = timeout
	}
	//进程收到退出或重启信号后，需要执行的方法
	if firstSweepFunc == nil {
		firstSweepFunc = func() error { return nil }
	}
	//退出与重启流程最大等待时间
	if beforeExitingFunc == nil {
		beforeExitingFunc = func() error { return nil }
	}
	that.firstSweep = func() error {
		return errors.Merge(
			firstSweepFunc(),    //执行自定义方法
			that.setInherited(), // 执行句柄继承方法
		)
	}
	that.beforeExiting = func() error {
		return errors.Merge(beforeExitingFunc())
	}
}

// IsChild 判断当前进程是在子进程还是父进程
func (that *graceful) isChild() bool {
	return that.childBool
}

// 保存监听的地址
func (that *graceful) setListenAddrMap(addr InheritAddr) {
	if that.listenAddrMap == nil {
		that.listenAddrMap = gmap.NewStrAnyMap(true)
	}
	// 做一次地址转换
	host := addr.Host
	if len(host) == 0 {
		host = "0.0.0.0"
	}
	if host == "0.0.0.0" {
		host = "[::]"
	}
	that.listenAddrMap.Set(fmt.Sprintf("%s:%d", host, addr.Port), addr)
}

// 初始化钩子函数
func (that *DServer) initGraceful() {
	normal.AddInheritedFunc(that.graceful.AddInherited)
	normal.GetInheritedFunc(that.graceful.GetInheritedFunc)
	quic.AddInheritedFunc(that.graceful.AddInheritedQUIC)
	quic.GetInheritedFunc(that.graceful.GetInheritedFuncQUIC)
	kcp.AddInheritedFunc(that.graceful.AddInheritedKCP)
	kcp.GetInheritedFunc(that.graceful.GetInheritedFuncKCP)
	that.graceful.setInherited = func() error {
		_ = normal.SetInherited() //把监听的文件句柄数量写入环境变量，方便子进程使用
		_ = quic.SetInherited()   // 把quic协议监听的文件句柄写入环境变量，方便子进程使用
		_ = kcp.SetInherited()    // 把kcp协议监听的文件句柄写入环境变量，方便子进程使用
		return nil
	}
}

// inheritListenerList 在多进程模式下，调用该方法，预先初始化监听
func (that *DServer) inheritListenerList() error {
	for _, addr := range that.inheritAddr {
		if addr.Network == "quic" || addr.Network == "kcp" {
			return gerror.Newf("Master-Worker进程模式不支持 quic,kcp协议")
		}
		network := that.translateNetwork(addr.Network)
		if addr.Network == "https" && addr.TlsConfig == nil {
			return gerror.Newf("使用https协议，必须传入证书")
		}
		err := that.inheritedListener(utils.NewFakeAddr(network, addr.Host, gconv.String(addr.Port)), addr.TlsConfig)
		if err != nil {
			return err
		}
		that.graceful.setListenAddrMap(addr)
		logger.Printf("多进程模式，主进程监听(network: %s,host: %s,port: %d)", addr.Network, addr.Host, addr.Port)
	}
	return nil
}

// inheritedListener 在多进程模式下，调用该方法，预先初始化监听
func (that *DServer) inheritedListener(addr net.Addr, tlsConfig *tls.Config) (err error) {

	if !that.isMaster() {
		return nil
	}
	addrStr := addr.String()
	network := addr.Network()
	var port string
	switch addrF := addr.(type) {
	case *utils.FakeAddr:
		_, port = addrF.Host(), addrF.Port()
	default:
		_, port, err = net.SplitHostPort(addrStr)
		if err != nil {
			return err
		}
	}
	if gstr.Contains(network, "tcp") && port == "0" {
		return gerror.New("必须明确的指定要监听的端口，不能使用随机端口")
	}

	lis, err := that.listen(network, addrStr)
	if err == nil && tlsConfig != nil {
		if len(tlsConfig.Certificates) == 0 && tlsConfig.GetCertificate == nil {
			return gerror.New("tls: neither Certificates nor GetCertificate set in Config")
		}
		lis = tls.NewListener(lis, tlsConfig)
	}
	if err != nil {
		return err
	}
	that.graceful.AddInherited([]net.Listener{lis}, nil)
	return nil
}

// 在多进程模式下，master进程预先监听地址
func (that *DServer) listen(nett, addr string) (net.Listener, error) {
	switch nett {
	default:
		return nil, net.UnknownNetworkError(nett)
	case "tcp", "tcp4", "tcp6":
		tcpAddr, err := net.ResolveTCPAddr(nett, addr)
		if err != nil {
			return nil, err
		}
		l, err := net.ListenTCP(nett, tcpAddr)
		if err != nil {
			return nil, err
		}
		return l, nil
	case "unix", "unixpacket", "invalid_unix_net_for_test":
		unixAddr, err := net.ResolveUnixAddr(nett, addr)
		if err != nil {
			return nil, err
		}
		l, err := net.ListenUnix(nett, unixAddr)
		if err != nil {
			return nil, err
		}
		return l, nil
	}
}

// 转换协议
func (that *DServer) translateNetwork(network string) string {
	switch network {
	case "tcp", "tcp4", "tcp6", "http", "https":
		return "tcp"
	case "unix", "unixpacket", "invalid_unix_net_for_test":
		return "unix"
	default:
		return "tcp"
	}
}
