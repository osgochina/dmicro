package graceful

import (
	"crypto/tls"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/errors/gerror"
	"github.com/osgochina/dmicro/utils/inherit"
	"net"
	"os"
	"os/exec"
	"time"
)

var defaultGraceful *graceful

func init() {
	defaultGraceful = newGraceful()
	inherit.AddInheritedFunc(defaultGraceful.AddInherited)
	defaultGraceful.SetShutdown(minShutdownTimeout, nil, nil)
}

//  进程收到结束或重启信号后，存活的最大时间
const minShutdownTimeout = 15 * time.Second

// GraceModel 平滑重启模型
type GraceModel int

const (
	// GraceChangeProcess 使用父子进程平滑重启
	GraceChangeProcess GraceModel = 1
	// GraceMasterWorker 使用master worker进程平滑重启
	GraceMasterWorker GraceModel = 2
)

const (
	// 当前是否是在子进程
	isChildKey = "GRACEFUL_IS_CHILD"
	// 父进程的监听列表
	parentAddrKey = "GRACEFUL_INHERIT_LISTEN_PARENT_ADDR"
	// 继承过来的fd有多少个
	envCountKey = "INHERIT_LISTEN_FDS"
)

//当前进程的状态
const (
	// 初始状态
	statusActionNone = 0
	// 进程在重启中
	statusActionRestarting = 1
	// 进程正在结束中
	statusActionShuttingDown = 2
)

type graceful struct {
	// 使用的模型
	model GraceModel
	// 当前进程的状态
	processStatus *gtype.Int

	//监听的信号
	signal chan os.Signal
	// 进程监听的地址列表，格式为 map[network][host]array[host:port]
	// 如：{"tcp":{"127.0.0.1":["127.0.0.1:8200","127.0.0.1:9980"]}}
	listenAddrList *gmap.StrAnyMap

	// 将要被子进程继承的环境变量
	inheritedEnv *gmap.StrStrMap
	// 将要被子进程继承的监听列表
	inheritedProcFiles *garray.Array

	// GracefulMasterWorker 模型使用，初始化需要监听的地址端口，方便子进程继承
	active []net.Listener

	// 进程收到退出或重启信号后，需要执行的方法
	firstSweep func() error
	// 进程真正退出前需要执行的方法
	beforeExiting func() error
	// 退出与重启流程最大等待时间
	shutdownTimeout time.Duration

	// 启动的endpoint列表
	endpointList *gset.Set

	shutdownCallback func(set *gset.Set) error

	// master worker模式下子进程命令句柄
	mwChildCmd chan *exec.Cmd
	// master worker模式下的子进程pid
	mwPid       int
	enableGHttp bool
}

type filer interface {
	File() (*os.File, error)
}

// Graceful 获取graceful对象
func Graceful() *graceful {
	return defaultGraceful
}

// InheritAddr master进程需要监听的配置
type InheritAddr struct {
	Network   string
	Host      string
	Port      string
	TlsConfig *tls.Config
}

// SetInheritListener 启动master worker模式的监听
func SetInheritListener(address []InheritAddr) error {
	defaultGraceful.SetModel(GraceMasterWorker)
	if !defaultGraceful.isChild() {
		var ch = make(chan int, 1)
		go func() {
			ch <- 1
			defaultGraceful.GraceSignal()
		}()
		<-ch
		for _, addr := range address {
			err := defaultGraceful.inheritedListener(inherit.NewFakeAddr(addr.Network, addr.Host, addr.Port), addr.TlsConfig)
			if err != nil {
				return err
			}
		}
		defaultGraceful.SetParentListenAddrList()
		cmd, err := defaultGraceful.startProcess()
		if err != nil {
			return gerror.Newf("启动子进程失败，error:%v", err)
		}
		defaultGraceful.mwChildCmd <- cmd
		defaultGraceful.MWWait()
	}
	return nil
}

// GraceSignal 监听信号
func GraceSignal() {
	defaultGraceful.onStart()
	defaultGraceful.GraceSignal()
}

// SetShutdown 设置重启钩子
func SetShutdown(timeout time.Duration, firstSweepFunc, beforeExitingFunc func() error) {
	defaultGraceful.SetShutdown(timeout, firstSweepFunc, beforeExitingFunc)
}

// Shutdown 停止服务
func Shutdown(timeout ...time.Duration) {
	defaultGraceful.Shutdown(timeout...)
}

// GetModel 当前进程模型
func GetModel() GraceModel {
	return defaultGraceful.model
}

// IsChild 判断当前是否是子进程
func IsChild() bool {
	return defaultGraceful.isChild()
}
