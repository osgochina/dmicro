package client

import (
	"context"
	"crypto/tls"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/selector"
	"time"
)

type Options struct {
	Context           context.Context // 上下文
	Network           string          // 网络类型
	LocalIP           string          // 本地网络
	TlsCertFile       string
	TlsKeyFile        string
	TLSConfig         *tls.Config
	ProtoFunc         proto.ProtoFunc
	SessionAge        time.Duration
	ContextAge        time.Duration
	DialTimeout       time.Duration
	SlowCometDuration time.Duration
	BodyCodec         string
	PrintDetail       bool
	CountTime         bool
	HeartbeatTime     time.Duration
	Registry          registry.Registry
	Selector          selector.Selector
	RetryTimes        int
	GlobalLeftPlugin  []drpc.Plugin
}

type Option func(*Options)

// NewOptions 初始化配置
func NewOptions(options ...Option) Options {
	opts := Options{
		Context:           context.Background(),
		Network:           "tcp",
		LocalIP:           "0.0.0.0",
		BodyCodec:         DefaultBodyCodec,
		SessionAge:        DefaultSessionAge,
		ContextAge:        DefaultContextAge,
		DialTimeout:       DefaultDialTimeout,
		SlowCometDuration: DefaultSlowCometDuration,
		RetryTimes:        DefaultRetryTimes,
		PrintDetail:       false,
		CountTime:         false,
		HeartbeatTime:     time.Duration(0),
		ProtoFunc:         drpc.DefaultProtoFunc(),
	}
	for _, o := range options {
		o(&opts)
	}

	return opts
}

func (that *Options) EndpointConfig() drpc.EndpointConfig {

	c := drpc.EndpointConfig{
		Network:           that.Network,
		LocalIP:           that.LocalIP,
		DefaultBodyCodec:  that.BodyCodec,
		DefaultSessionAge: that.SessionAge,
		DefaultContextAge: that.ContextAge,
		SlowCometDuration: that.SlowCometDuration,
		PrintDetail:       that.PrintDetail,
		CountTime:         that.CountTime,
		DialTimeout:       that.DialTimeout,
	}
	return c
}

// Registry 设置服务注册中心
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
		// set in the selector
		_ = o.Selector.Init(selector.Registry(r))
	}
}

// Selector 设置选择器
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

// GlobalLeftPlugin 设置插件
func GlobalLeftPlugin(plugin ...drpc.Plugin) Option {
	return func(o *Options) {
		o.GlobalLeftPlugin = append(o.GlobalLeftPlugin, plugin...)
	}
}

// HeartbeatTime 设置心跳包时间
func HeartbeatTime(t time.Duration) Option {
	return func(o *Options) {
		o.HeartbeatTime = t
	}
}

// TlsFile 设置证书内容
func TlsFile(tlsCertFile string, tlsKeyFile string) Option {
	return func(o *Options) {
		o.TlsCertFile = tlsCertFile
		o.TlsKeyFile = tlsKeyFile
	}
}

// TlsConfig 设置证书对象
func TlsConfig(config *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = config
	}
}

// ProtoFunc 设置协议方法
func ProtoFunc(pf proto.ProtoFunc) Option {
	return func(o *Options) {
		o.ProtoFunc = pf
	}
}

// RetryTimes 设置重试次数
func RetryTimes(n int) Option {
	return func(o *Options) {
		o.RetryTimes = n
	}
}
func SessionAge(n time.Duration) Option {
	return func(o *Options) {
		o.SessionAge = n
	}
}

func ContextAge(n time.Duration) Option {
	return func(o *Options) {
		o.ContextAge = n
	}
}

func DialTimeout(n time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = n
	}
}

func SlowCometDuration(n time.Duration) Option {
	return func(o *Options) {
		o.SlowCometDuration = n
	}
}

func BodyCodec(c string) Option {
	return func(o *Options) {
		o.BodyCodec = c
	}
}
func PrintDetail(c bool) Option {
	return func(o *Options) {
		o.PrintDetail = c
	}
}
func CountTime(c bool) Option {
	return func(o *Options) {
		o.CountTime = c
	}
}
func Network(net string) Option {
	return func(o *Options) {
		o.Network = net
	}
}
func LocalIP(addr string) Option {
	return func(o *Options) {
		o.LocalIP = addr
	}
}
