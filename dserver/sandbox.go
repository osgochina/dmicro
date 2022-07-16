package dserver

import (
	"context"
	"github.com/gogf/gf/os/gtime"
	"github.com/osgochina/dmicro/supervisor/process"
)

// ISandbox 服务沙盒的接口
type ISandbox interface {
	Name() string    // 沙盒名字
	Setup() error    // 启动沙盒
	Shutdown() error //关闭沙盒
}

// BaseSandbox sandbox的基类，必须继承它
type BaseSandbox struct {
	Service *DService
	Config  *Config
	Context context.Context
}

// sandbox的容器
type sandboxContainer struct {
	sandbox  ISandbox
	started  *gtime.Time   //服务启动时间
	stopTime *gtime.Time   //服务关闭时间
	state    process.State // sandbox的运行状态
}
