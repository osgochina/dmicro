package registry

import (
	"context"
	"crypto/tls"
	"time"
)

// Options 配置
type Options struct {
	AddrList  []string        // 地址列表
	Timeout   time.Duration   // 超时时间
	Secure    bool            // 是否加密
	TLSConfig *tls.Config     // tls证书配置
	Context   context.Context // 上下文系信息
}

// RegisterOptions 注册配置信息
type RegisterOptions struct {
	TTL     time.Duration
	Context context.Context
}

// WatchOptions 监视器参数
type WatchOptions struct {
	Service string
	Context context.Context
}

// DeregisterOptions 取消注册参数
type DeregisterOptions struct {
	Context context.Context
}

type GetOptions struct {
	Context context.Context
}

type ListOptions struct {
	Context context.Context
}

// OptAddrList 设置地址列表
func OptAddrList(addrList ...string) Option {
	return func(o *Options) {
		o.AddrList = addrList
	}
}

// OptTimeout 设置超时时间
func OptTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// OptSecure 是否加密
func OptSecure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// OptTLSConfig 设置tls证书信息
func OptTLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// OptServiceName 设置服务名称
func OptServiceName(name string) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "ServiceName", name)
	}
}

// OptServiceVersion 设置服务版本
func OptServiceVersion(version string) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "ServiceVersion", version)
	}
}

type leasesInterval struct{}

// OptLeasesInterval 租约续期时间
func OptLeasesInterval(t time.Duration) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, leasesInterval{}, t)
	}
}

// OptRegisterTTL 设置服务的生存时间
func OptRegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

// OptRegisterContext 设置注册服务的上下文
func OptRegisterContext(ctx context.Context) RegisterOption {
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

// OptWatchService 监视器监听指定的服务
func OptWatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}

// OptWatchContext 监视器的上下文
func OptWatchContext(ctx context.Context) WatchOption {
	return func(o *WatchOptions) {
		o.Context = ctx
	}
}

func OptDeregisterContext(ctx context.Context) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Context = ctx
	}
}

func OptGetContext(ctx context.Context) GetOption {
	return func(o *GetOptions) {
		o.Context = ctx
	}
}

func OptListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}
