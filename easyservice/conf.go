package easyservice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"net"
	"strconv"
	"time"
)

//Deprecated
type BoxConf struct {
	id                 int           //服务沙盒的id"
	SandBoxName        string        `json:"sandbox_name"   comment:"服务沙盒的名字"`
	Network            string        `json:"network"        comment:"使用的网络协议; tcp, tcp4, tcp6,kcp,quic,unix or unixpacket"`
	ListenAddress      string        `json:"listen_address" comment:"监听地址"`
	PrintDetail        bool          `json:"print_detail"   comment:"是否显示通讯详情"`
	SessionMaxTimeout  time.Duration `json:"session_max_timeout" comment:"会话生命周期"`
	ResponseMaxTimeout time.Duration `json:"response_max_timeout" comment:"单次处理响应最长超时时间"`
	SlowTimeout        time.Duration `json:"slow_timeout" comment:"慢请求时间"`
	RequestMaxTimeout  time.Duration `json:"request_max_timeout" comment:"请求最长超时时间"`
	RedialTimes        int           `json:"redial_times" comment:"在链接中断时候，试图链接服务端的最大重试次数。仅限客户端角色使用"`
	RedialInterval     time.Duration `json:"redial_interval" comment:"仅限客户端角色使用 试图链接服务端时候，重试的时间间隔"`
	ctx                context.Context
}

// NewBoxConf 创建BoxConf对象
//Deprecated
func NewBoxConf(name string, config *gcfg.Config, parsers ...*gcmd.Parser) *BoxConf {
	cfg := &BoxConf{}
	cfg.ctx = context.TODO()
	cfg.id = config.MustGet(cfg.ctx, fmt.Sprintf("%s.Id", name), 0).Int()
	cfg.SandBoxName = config.MustGet(cfg.ctx, fmt.Sprintf("%s.Name", name), "").String()

	cfg.Network = config.MustGet(cfg.ctx, fmt.Sprintf("%s.Network", name), "tcp").String()
	host := config.MustGet(cfg.ctx, fmt.Sprintf("%s.Host", name), "0.0.0.0").String()
	port := config.MustGet(cfg.ctx, fmt.Sprintf("%s.Port", name), 0).Int()
	cfg.ListenAddress = fmt.Sprintf("%s:%d", host, port)

	cfg.SessionMaxTimeout = config.MustGet(cfg.ctx, fmt.Sprintf("%s.SessionMaxTimeout", name), 0).Duration()
	cfg.ResponseMaxTimeout = config.MustGet(cfg.ctx, fmt.Sprintf("%s.ResponseMaxTimeout", name), 0).Duration()
	cfg.RequestMaxTimeout = config.MustGet(cfg.ctx, fmt.Sprintf("%s.RequestMaxTimeout", name), 0).Duration()
	cfg.SlowTimeout = config.MustGet(cfg.ctx, fmt.Sprintf("%s.SlowTimeout", name), 0).Duration()
	cfg.RedialTimes = config.MustGet(cfg.ctx, fmt.Sprintf("%s.RedialTimes", name), 0).Int()
	cfg.RedialInterval = config.MustGet(cfg.ctx, fmt.Sprintf("%s.RedialTimes", name), 0).Duration()

	debug := config.MustGet(cfg.ctx, "Debug", false).Bool()
	cfg.PrintDetail = debug
	// 命令行参数覆盖配置文件
	if len(parsers) > 0 {
		parser := parsers[0]
		cfg.ListenAddress = fmt.Sprintf("%s:%d",
			parser.GetOpt("host", host).String(),
			parser.GetOpt("port", port).Int(),
		)
		cfg.Network = parser.GetOpt(fmt.Sprintf("%s.Network", name), cfg.Network).String()
	}

	return cfg
}

// DefaultBoxConf 创建运行配置
//Deprecated
func DefaultBoxConf(parser *gcmd.Parser, config *gcfg.Config) *BoxConf {
	cfg := NewBoxConf("default.sandbox", config, parser)
	return cfg
}

// ListenPort 获取监听端口
func (that *BoxConf) ListenPort() string {
	_, port, err := net.SplitHostPort(that.ListenAddress)
	if err != nil {
		logger.Fatalf(that.ctx, "%v", err)
	}
	return port
}

// InnerIpPort 获取内网的服务地址
func (that *BoxConf) InnerIpPort() string {
	host, err := gipv4.GetIntranetIp()
	if err != nil {
		logger.Fatalf(that.ctx, "%v", err)
	}
	return fmt.Sprintf("%s:%s", host, that.ListenPort())
}

func (that *BoxConf) SetId(id int) {
	that.id = id
}

func (that *BoxConf) GetId() int {
	return that.id
}

// EndpointConfig 返回rpc服务要用的配置文件
//Deprecated
func (that *BoxConf) EndpointConfig() drpc.EndpointConfig {
	var c = drpc.EndpointConfig{
		PrintDetail:       that.PrintDetail,
		Network:           that.Network,
		DefaultSessionAge: that.SessionMaxTimeout,
		DefaultContextAge: that.ResponseMaxTimeout,
		SlowCometDuration: that.SlowTimeout,
	}
	if len(that.ListenAddress) > 0 {
		host, port, err := net.SplitHostPort(that.ListenAddress)
		if err != nil {
			logger.Fatalf(that.ctx, "%v", err)
		}
		listenPort, _ := strconv.Atoi(port)
		c.LocalIP = host
		c.ListenPort = uint16(listenPort)
	}
	return c
}
