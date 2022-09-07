package dserver

import (
	"context"
	"github.com/gogf/gf/os/gtime"
	"github.com/osgochina/dmicro/supervisor/process"
)

// KindSandbox sandbox种类
type kindSandbox string

const (
	serviceKindSandbox kindSandbox = "service" // 持续提供服务
)

// ISandbox 服务沙盒的接口
type ISandbox interface {
	Name() string    // 沙盒名字
	Setup() error    // 启动沙盒
	Shutdown() error //关闭沙盒
}

// BaseSandbox sandbox的基类，必须继承它
type BaseSandbox struct {
	Context context.Context
	Service *DService
	Config  *Config
}

// ServiceSandbox 服务
type ServiceSandbox struct {
	BaseSandbox
}

// sandbox的容器
type sandboxContainer struct {
	sandbox  ISandbox
	kind     kindSandbox   // 容器的种类
	started  *gtime.Time   //服务启动时间
	stopTime *gtime.Time   //服务关闭时间
	state    process.State // sandbox的运行状态
}
