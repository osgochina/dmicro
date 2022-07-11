package dserver

import "context"

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
