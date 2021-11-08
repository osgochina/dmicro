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
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/errors"
	"github.com/osgochina/dmicro/utils/inherit"
	"net"
	"os"
	"strconv"
	"time"
)

// newGraceful 创建对象
func newGraceful() *graceful {
	return &graceful{
		model:              GraceChangeProcess,
		processStatus:      gtype.NewInt(statusActionNone),
		signal:             make(chan os.Signal),
		listenAddrList:     gmap.NewStrAnyMap(true),
		inheritedEnv:       gmap.NewStrStrMap(true),
		inheritedProcFiles: garray.NewArray(true),
		active:             make([]net.Listener, 0),
		firstSweep:         func() error { return nil },
		beforeExiting:      func() error { return nil },
		endpointList:       gset.New(true),
	}
}

// SetModel 设置模式
func (that *graceful) SetModel(model GraceModel) {
	that.model = model
}

// IsChild 判断当前进程是在子进程还是父进程
func (that *graceful) IsChild() bool {
	isWorker := genv.GetVar(isChildKey, nil)
	if isWorker.IsNil() {
		return false
	}
	if isWorker.Bool() == true {
		return true
	}
	return false
}

// SetParentListenAddrList 设置已监听的地址列表到环境变量，在子进程启动的时候，把该环境变量注入到启动参数中
// 在父进程收到平滑重启信号以后，会调用该方法
func (that *graceful) SetParentListenAddrList() {
	env := make(map[string]string)
	env[isChildKey] = "true"
	j, err := that.listenAddrList.MarshalJSON()
	if err != nil {
		logger.Error(err)
	} else {
		env[parentAddrKey] = gconv.String(j)
	}
	var procFiles []*os.File
	// master-worker进程模型逻辑
	if that.model == GraceMasterWorker && len(that.active) > 0 {
		procFiles = make([]*os.File, 0, len(that.active))
		for _, l := range that.active {
			f, e := l.(filer).File()
			if e != nil {
				logger.Error(e)
				continue
			}
			procFiles = append(procFiles, f)
		}
		env[envCountKey] = strconv.Itoa(len(procFiles))
	}

	that.AddInherited(procFiles, env)
}

// InitParentAddrList 通过环境变量，初始化父进程监听的端口
// 在服务启动的时候，首先从环境变量中获取父进程监听的地址端口，
// 如果是首次启动或者当前进程是父进程，则不会获取到这些数据
// 如果是优雅的无缝重启的子进程，则能通过环境变量获取到这些数据，从而复用链接，做到无缝重启
func (that *graceful) InitParentAddrList() {
	parentAddr := genv.GetVar(parentAddrKey, nil)
	if parentAddr.IsNil() {
		return
	}
	json := gjson.New(parentAddr)
	if json.IsNil() {
		return
	}
	err := that.listenAddrList.UnmarshalValue(json)
	if err != nil {
		logger.Error(err)
	}
}

// PushParentAddr 把监听的地址端口写入到变量，优雅重启的时候写入到环境变量，让子进程使用
// listenAddrList变量的格式是 gmap.StrAnyMap(network,gmap.StrAnyMap(host,garray.StrArray(addr)))
func (that *graceful) PushParentAddr(network, host, addr string) {
	that.unifyLocalhost(&host)
	nw, found := that.listenAddrList.Search(network)
	if !found {
		nw = gmap.NewStrAnyMap(true)
		that.listenAddrList.Set(network, nw)
	}
	if nwMap, ok := nw.(*gmap.StrAnyMap); ok {
		hs, f := nwMap.Search(host)
		if !f {
			hs = garray.NewStrArray(true)
			nwMap.Set(host, hs)
		}
		if ar, ok := hs.(*garray.StrArray); ok {
			ar.Append(addr)
		}
	}
}

// PopParentAddr 从监听变量中出栈指定的地址端口
func (that *graceful) PopParentAddr(network, host, addr string) string {
	that.unifyLocalhost(&host)
	nw, found := that.listenAddrList.Search(network)
	if !found {
		return addr
	}
	nwMap, ok := nw.(*gmap.StrAnyMap)
	if !ok {
		return addr
	}
	hs, ok := nwMap.Search(host)
	if !ok {
		return addr
	}
	ar, ok := hs.(*garray.StrArray)
	if !ok {
		return addr
	}
	if ar.Len() == 0 {
		return addr
	}
	a, f := ar.PopLeft()
	if !f {
		return addr
	}
	return a
}

// 针对地址格式做统一的转换
func (that *graceful) unifyLocalhost(host *string) {
	switch *host {
	case "localhost":
		*host = "127.0.0.1"
	case "0.0.0.0":
		*host = "::"
	}
}

// inheritedListener 在GracefulMasterWorker模式下，调用该方法，预先初始化监听
func (that *graceful) inheritedListener(addr net.Addr, tlsConfig *tls.Config) (err error) {
	if that.model != GraceMasterWorker {
		return gerror.New("必须为GracefulMasterWorker模式才可以调用InheritedListener方法")
	}
	if that.IsChild() {
		return nil
	}
	addrStr := addr.String()
	network := addr.Network()
	var host, port string
	switch addrF := addr.(type) {
	case *inherit.FakeAddr:
		host, port = addrF.Host(), addrF.Port()
	default:
		host, port, err = net.SplitHostPort(addrStr)
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
	if err == nil {
		that.PushParentAddr(network, host, lis.Addr().String())
	}
	that.active = append(that.active, lis)
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
		defaultGraceful.SetParentListenAddrList()
		return errors.Merge(
			firstSweepFunc(),       //执行自定义方法
			inherit.SetInherited(), //把监听的文件句柄数量写入环境变量，方便子进程使用
		)
	}
	that.beforeExiting = func() error {
		return errors.Merge(defaultGraceful.shutdownEndpoint(), beforeExitingFunc())
	}
}

// MWWait master worker 模式的主进程等待子进程运行
func (that *graceful) MWWait() {
	for {
		var err error
		that.mwChildCmd, err = that.startProcess()
		if err != nil {
			logger.Errorf("启动子进程失败，error:%v", err)
			return
		}
		logger.Infof("Master-Worker模式启动子进程成功，父进程:%d子进程:%d", os.Getpid(), that.mwChildCmd.Process.Pid)
		err = that.mwChildCmd.Wait()
		if err != nil {
			logger.Warningf("子进程:%d 非正常退出，退出原因:%v", that.mwChildCmd.Process.Pid, err)
		} else {
			logger.Infof("子进程:%d 正常退出", that.mwChildCmd.Process.Pid)
		}
	}
}
