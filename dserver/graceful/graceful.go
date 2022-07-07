package graceful

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gtype"
	"github.com/osgochina/dmicro/drpc/netproto/kcp"
	"github.com/osgochina/dmicro/drpc/netproto/normal"
	"github.com/osgochina/dmicro/drpc/netproto/quic"
	"os"
	"time"
)

var defaultGraceful *graceful

func init() {
	defaultGraceful = &graceful{
		processStatus:         gtype.NewInt(statusActionNone),
		signal:                make(chan os.Signal),
		inheritedEnv:          gmap.NewStrStrMap(true),
		inheritedProcListener: garray.NewArray(true),
		firstSweep:            func() error { return nil },
		beforeExiting:         func() error { return nil },
	}
	normal.AddInheritedFunc(AddInherited)
	normal.GetInheritedFunc(GetInheritedFunc)
	quic.AddInheritedFunc(AddInheritedQUIC)
	quic.GetInheritedFunc(GetInheritedFuncQUIC)
	kcp.AddInheritedFunc(AddInheritedKCP)
	kcp.GetInheritedFunc(GetInheritedFuncKCP)
	SetInherited(func() error {
		_ = normal.SetInherited() //把监听的文件句柄数量写入环境变量，方便子进程使用
		_ = quic.SetInherited()   // 把quic协议监听的文件句柄写入环境变量，方便子进程使用
		_ = kcp.SetInherited()    // 把kcp协议监听的文件句柄写入环境变量，方便子进程使用
		return nil
	})
	defaultGraceful.SetShutdown(minShutdownTimeout, nil, nil)

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

//  进程收到结束或重启信号后，存活的最大时间
const minShutdownTimeout = 15 * time.Second

const (
	// 当前是否是在子进程
	isChildKey = "GRACEFUL_IS_CHILD"
	// 父进程的监听列表
	parentAddrKey = "GRACEFUL_INHERIT_LISTEN_PARENT_ADDR"
	// 继承过来的fd有多少个
	//envCountKey = "INHERIT_LISTEN_FDS"
	// gf框架的ghttp服务平滑重启key
	adminActionReloadEnvKey = "GF_SERVER_RELOAD"
)

type graceful struct {
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
}

// GraceSignal 监听信号
func GraceSignal(model int) {
	defaultGraceful.onStart()
	defaultGraceful.GraceSignal(model)
}

// SetShutdown 设置重启钩子
func SetShutdown(timeout time.Duration, firstSweepFunc, beforeExitingFunc func() error) {
	defaultGraceful.SetShutdown(timeout, firstSweepFunc, beforeExitingFunc)
}

// Shutdown 停止服务
func Shutdown(timeout ...time.Duration) {
	defaultGraceful.shutdown(timeout...)
}
