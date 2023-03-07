package server

import (
	"context"
	"crypto/tls"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/metrics"
	"github.com/osgochina/dmicro/registry"
	"net"
	"time"
)

type Options struct {
	Context           context.Context // 上下文
	ServiceName       string          // 服务名称
	ServiceVersion    string          // 服务版本
	Network           string          // 网络类型
	ListenAddress     string          // 要监听的地址
	TlsCertFile       string
	TlsKeyFile        string
	TLSConfig         *tls.Config
	ProtoFunc         proto.ProtoFunc
	SessionAge        time.Duration
	ContextAge        time.Duration
	SlowCometDuration time.Duration
	BodyCodec         string
	PrintDetail       bool
	Registry          registry.Registry
	GlobalPlugin      []drpc.Plugin
	EnableHeartbeat   bool
	Metrics           metrics.Metrics // 统计信息
}

type Option func(*Options)

// newOptions 初始化配置
func newOptions(options ...Option) Options {
	opts := Options{
		Context:           context.Background(),
		Network:           "tcp",
		ListenAddress:     "0.0.0.0:0",
		BodyCodec:         defaultBodyCodec,
		SessionAge:        defaultSessionAge,
		ContextAge:        defaultContextAge,
		SlowCometDuration: defaultSlowCometDuration,
		PrintDetail:       false,
		EnableHeartbeat:   false,
		ProtoFunc:         drpc.DefaultProtoFunc(),
	}
	for _, o := range options {
		o(&opts)
	}

	return opts
}

func (that Options) EndpointConfig() drpc.EndpointConfig {

	c := drpc.EndpointConfig{
		Network:           that.Network,
		DefaultBodyCodec:  that.BodyCodec,
		DefaultSessionAge: that.SessionAge,
		DefaultContextAge: that.ContextAge,
		SlowCometDuration: that.SlowCometDuration,
		PrintDetail:       that.PrintDetail,
	}
	switch that.Network {
	case "tcp", "tcp4", "tcp6", "kcp", "udp", "udp4", "udp6", "quic":
		ip, port, err := net.SplitHostPort(that.ListenAddress)
		if err != nil {
			logger.Fatal(context.TODO(), err)
		}
		c.ListenIP = ip
		c.ListenPort = gconv.Uint16(port)
	case "unix", "unixpacket":
		c.ListenIP = that.ListenAddress
		c.ListenPort = 0
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
	}
}

// OptGlobalPlugin 设置插件
func OptGlobalPlugin(plugin ...drpc.Plugin) Option {
	return func(o *Options) {
		o.GlobalPlugin = append(o.GlobalPlugin, plugin...)
	}
}

// OptEnableHeartbeat 是否开启心跳包
func OptEnableHeartbeat(t bool) Option {
	return func(o *Options) {
		o.EnableHeartbeat = t
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

// OptProtoFunc 设置协议
func OptProtoFunc(pf proto.ProtoFunc) Option {
	return func(o *Options) {
		o.ProtoFunc = pf
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

// OptListenAddress 设置监听的网络地址
func OptListenAddress(addr string) Option {
	return func(o *Options) {
		o.ListenAddress = addr
	}
}

// OptMetrics 设置统计数据对象
func OptMetrics(m metrics.Metrics) Option {
	return func(o *Options) {
		o.Metrics = m
	}
}
