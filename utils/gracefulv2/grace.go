package gracefulv2

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/container/gtype"
	"github.com/osgochina/dmicro/utils/errors"
	"github.com/osgochina/dmicro/utils/inherit"
	"net"
	"os"
	"os/exec"
	"time"
)

var defaultGraceful *Graceful

func init() {
	defaultGraceful = newGraceful(GracefulNormal)
	inherit.AddInheritedFunc(defaultGraceful.AddInherited)
	defaultGraceful.SetShutdown(minShutdownTimeout, nil, nil)
}

//  进程收到结束或重启信号后，存活的最大时间
const minShutdownTimeout = 15 * time.Second

// GracefulModel 平滑重启模型
type GracefulModel int

const (
	// GracefulNormal 不适用平滑重启
	GracefulNormal GracefulModel = 0
	// GracefulChangeProcess 使用父子进程平滑重启
	GracefulChangeProcess GracefulModel = 1
	// GracefulMasterWorker 使用master worker进程平滑重启
	GracefulMasterWorker GracefulModel = 2
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

type Graceful struct {
	// 使用的模型
	model GracefulModel
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

	// master worker模式下的子进程pid
	masterWorkerChildCmd *exec.Cmd
}

type filer interface {
	File() (*os.File, error)
}

// GetGraceful 获取graceful对象
func GetGraceful() *Graceful {
	return defaultGraceful
}

// SetShutdown 设置退出的基本参数
// timeout 传入负数，表示永远不过期
func (that *Graceful) SetShutdown(timeout time.Duration, firstSweepFunc, beforeExitingFunc func() error) {
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
