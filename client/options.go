package client

import (
	"context"
	"crypto/tls"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/metrics"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/registry/memory"
	"github.com/osgochina/dmicro/selector"
	"time"
)

type Options struct {
	Context           context.Context // 上下文
	ServiceName       string          // 服务名称
	ServiceVersion    string          // 服务版本
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
	HeartbeatTime     time.Duration
	RetryTimes        int
	GlobalPlugin      []drpc.Plugin
	Registry          registry.Registry
	Selector          selector.Selector
	Metrics           metrics.Metrics // 统计信息
}

type Option func(*Options)

// NewOptions 初始化配置
func NewOptions(options ...Option) Options {
	opts := Options{
		Context:           context.Background(),
		Network:           "tcp",
		LocalIP:           "0.0.0.0",
		BodyCodec:         defaultBodyCodec,
		SessionAge:        defaultSessionAge,
		ContextAge:        defaultContextAge,
		DialTimeout:       defaultDialTimeout,
		SlowCometDuration: defaultSlowCometDuration,
		RetryTimes:        defaultRetryTimes,
		PrintDetail:       false,
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
		DialTimeout:       that.DialTimeout,
	}
	return c
}

// OptServiceName 设置服务名称
func OptServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

// OptServiceVersion 当前服务版本
func OptServiceVersion(version string) Option {
	return func(o *Options) {
		o.ServiceVersion = version
	}
}

// OptRegistry 设置服务注册中心
func OptRegistry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
		// 初始化默认selector
		if o.Selector == nil {
			o.Selector = selector.NewSelector(selector.OptRegistry(r))
		} else {
			_ = o.Selector.Init(selector.OptRegistry(r))
		}
	}
}

// OptSelector 设置选择器
func OptSelector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

// OptGlobalPlugin 设置插件
func OptGlobalPlugin(plugin ...drpc.Plugin) Option {
	return func(o *Options) {
		o.GlobalPlugin = append(o.GlobalPlugin, plugin...)
	}
}

// OptHeartbeatTime 设置心跳包时间
func OptHeartbeatTime(t time.Duration) Option {
	return func(o *Options) {
		o.HeartbeatTime = t
	}
}

// OptTlsFile 设置证书内容
func OptTlsFile(tlsCertFile string, tlsKeyFile string) Option {
	return func(o *Options) {
		o.TlsCertFile = tlsCertFile
		o.TlsKeyFile = tlsKeyFile
	}
}

// OptTlsConfig 设置证书对象
func OptTlsConfig(config *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = config
	}
}

// OptProtoFunc 设置协议方法
func OptProtoFunc(pf proto.ProtoFunc) Option {
	return func(o *Options) {
		o.ProtoFunc = pf
	}
}

// OptRetryTimes 设置重试次数
func OptRetryTimes(n int) Option {
	return func(o *Options) {
		o.RetryTimes = n
	}
}

// OptSessionAge 设置会话生命周期
func OptSessionAge(n time.Duration) Option {
	return func(o *Options) {
		o.SessionAge = n
	}
}

// OptContextAge 设置单次请求生命周期
func OptContextAge(n time.Duration) Option {
	return func(o *Options) {
		o.ContextAge = n
	}
}

// OptSlowCometDuration 设置慢请求的定义时间
func OptSlowCometDuration(n time.Duration) Option {
	return func(o *Options) {
		o.SlowCometDuration = n
	}
}

// OptBodyCodec 设置消息内容编解码器
func OptBodyCodec(c string) Option {
	return func(o *Options) {
		o.BodyCodec = c
	}
}

// OptPrintDetail 是否打印消息详情
func OptPrintDetail(c bool) Option {
	return func(o *Options) {
		o.PrintDetail = c
	}
}

// OptNetwork 设置网络类型
func OptNetwork(net string) Option {
	return func(o *Options) {
		o.Network = net
	}
}

// OptLocalIP 设置本地监听的地址
func OptLocalIP(addr string) Option {
	return func(o *Options) {
		o.LocalIP = addr
	}
}

// OptCustomService 自定义service
func OptCustomService(service *registry.Service) Option {
	return func(o *Options) {
		o.ServiceVersion = "1.0.0"
		if len(service.Name) > 0 {
			o.ServiceName = service.Name
		} else {
			service.Name = o.ServiceName
		}
		if len(service.Version) > 0 {
			o.ServiceVersion = service.Version
		} else {
			service.Version = o.ServiceVersion
		}
		o.Registry = memory.NewRegistry()
		err := o.Registry.Register(service)
		if err != nil {
			logger.Fatal(context.TODO(), err)
		}
		// 初始化默认selector
		if o.Selector == nil {
			o.Selector = selector.NewSelector(selector.OptRegistry(o.Registry))
		} else {
			err = o.Selector.Init(selector.OptRegistry(o.Registry))
			if err != nil {
				logger.Fatal(context.TODO(), err)
			}
		}
	}
}

// OptMetrics 设置统计数据接口
func OptMetrics(m metrics.Metrics) Option {
	return func(o *Options) {
		o.Metrics = m
	}
}
