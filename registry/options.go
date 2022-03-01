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

// AddrList 设置地址列表
func AddrList(addrList ...string) Option {
	return func(o *Options) {
		o.AddrList = addrList
	}
}

// Timeout 设置超时时间
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// Secure 是否加密
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// TLSConfig 设置tls证书信息
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// ServiceName 设置服务名称
func ServiceName(name string) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "ServiceName", name)
	}
}

// ServiceVersion 设置服务版本
func ServiceVersion(version string) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "ServiceVersion", version)
	}
}

type leasesInterval struct{}

// LeasesInterval 租约续期时间
func LeasesInterval(t time.Duration) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, leasesInterval{}, t)
	}
}

// RegisterTTL 设置服务的生存时间
func RegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

// RegisterContext 设置注册服务的上下文
func RegisterContext(ctx context.Context) RegisterOption {
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

// WatchService 监视器监听指定的服务
func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}

// WatchContext 监视器的上下文
func WatchContext(ctx context.Context) WatchOption {
	return func(o *WatchOptions) {
		o.Context = ctx
	}
}

func DeregisterContext(ctx context.Context) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Context = ctx
	}
}

func GetContext(ctx context.Context) GetOption {
	return func(o *GetOptions) {
		o.Context = ctx
	}
}

func ListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}
