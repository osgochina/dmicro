package graceful

import (
	"crypto/tls"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/utils/errors"
	"github.com/osgochina/dmicro/utils/inherit"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// newGraceful 创建对象
func newGraceful() *graceful {
	return &graceful{
		model:                 GraceChangeProcess,
		processStatus:         gtype.NewInt(statusActionNone),
		signal:                make(chan os.Signal),
		inheritedEnv:          gmap.NewStrStrMap(true),
		inheritedProcListener: garray.NewArray(true),
		firstSweep:            func() error { return nil },
		beforeExiting:         func() error { return nil },
		endpointList:          gset.New(true),
		mwChildCmd:            make(chan *exec.Cmd, 1),
	}
}

// SetModel 设置模式
func (that *graceful) SetModel(model GraceModel) {
	that.model = model
}

// isChild 判断当前进程是在子进程还是父进程
func (that *graceful) isChild() bool {
	isWorker := genv.GetVar(isChildKey, nil)
	if isWorker.IsNil() {
		return false
	}
	if isWorker.Bool() == true {
		return true
	}
	return false
}

// GetInheritedFunc 获取继承的fd
func (that *graceful) GetInheritedFunc() []int {
	parentAddr := genv.GetVar(parentAddrKey, nil)
	if parentAddr.IsNil() {
		return nil
	}
	json := gjson.New(parentAddr)
	if json.IsNil() {
		return nil
	}
	var fds []int
	for _, v := range json.Map() {
		fdv := gconv.String(v)
		if len(fdv) > 0 {
			for _, item := range gstr.SplitAndTrim(fdv, ",") {
				array := strings.Split(item, "#")
				fd := gconv.Int(array[1])
				if fd > 0 {
					fds = append(fds, fd)
				}
			}
		}
	}
	return fds
}

// inheritedListener 在GracefulMasterWorker模式下，调用该方法，预先初始化监听
func (that *graceful) inheritedListener(addr net.Addr, tlsConfig *tls.Config) (err error) {
	if that.model != GraceMasterWorker {
		return gerror.New("必须为GracefulMasterWorker模式才可以调用InheritedListener方法")
	}
	if that.isChild() {
		return nil
	}
	addrStr := addr.String()
	network := addr.Network()
	var port string
	switch addrF := addr.(type) {
	case *inherit.FakeAddr:
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
	that.inheritedProcListener.Append(lis)
	return nil
}

// 在GracefulMasterWorker模式下，master进程预先监听地址
func (that *graceful) listen(nett, addr string) (net.Listener, error) {
	if that.model != GraceMasterWorker {
		return nil, gerror.New("必须为GracefulMasterWorker模式才可以调用listen方法")
	}
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

// SetShutdown 设置退出的基本参数
// timeout 传入负数，表示永远不过期
func (that *graceful) SetShutdown(timeout time.Duration, firstSweepFunc, beforeExitingFunc func() error) {
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
			firstSweepFunc(),       //执行自定义方法
			inherit.SetInherited(), //把监听的文件句柄数量写入环境变量，方便子进程使用
		)
	}
	that.beforeExiting = func() error {
		return errors.Merge(beforeExitingFunc(), defaultGraceful.shutdownEndpoint())
	}
}
