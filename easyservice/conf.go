package easyservice

import (
	"fmt"
	"github.com/gogf/gf/net/gipv4"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gcmd"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"net"
	"strconv"
	"time"
)

type BoxConf struct {
	id                 int           //服务沙盒的id"
	SandBoxName        string        `json:"sandbox_name"   comment:"服务沙盒的名字"`
	Network            string        `json:"network"        comment:"使用的网络协议; tcp, tcp4, tcp6,kcp,quic,unix or unixpacket"`
	ListenAddress      string        `json:"listen_address" comment:"监听地址"`
	PrintDetail        bool          `json:"print_detail"   comment:"是否显示通讯详情"`
	SessionMaxTimeout  time.Duration `json:"session_max_timeout" comment:"会话生命周期"`
	ResponseMaxTimeout time.Duration `json:"response_max_timeout" comment:"单次处理响应最长超时时间"`
	SlowTimeout        time.Duration `json:"slow_timeout" comment:"慢请求时间"`
	CountTime          bool          `json:"count_time" comment:"是否统计消耗时间"`
	RequestMaxTimeout  time.Duration `json:"request_max_timeout" comment:"请求最长超时时间"`
	RedialTimes        int           `json:"redial_times" comment:"在链接中断时候，试图链接服务端的最大重试次数。仅限客户端角色使用"`
	RedialInterval     time.Duration `json:"redial_interval" comment:"仅限客户端角色使用 试图链接服务端时候，重试的时间间隔"`
}

// NewBoxConf 创建BoxConf对象
func NewBoxConf(name string, config *gcfg.Config, parsers ...*gcmd.Parser) *BoxConf {
	cfg := &BoxConf{}

	cfg.id = config.GetInt(fmt.Sprintf("%s.Id", name), 0)
	cfg.SandBoxName = config.GetString(fmt.Sprintf("%s.Name", name), "")

	cfg.Network = config.GetString(fmt.Sprintf("%s.Network", name), "tcp")
	host := config.GetString(fmt.Sprintf("%s.Host", name), "0.0.0.0")
	port := config.GetInt(fmt.Sprintf("%s.Port", name), 0)
	cfg.ListenAddress = fmt.Sprintf("%s:%d", host, port)

	cfg.SessionMaxTimeout = config.GetDuration(fmt.Sprintf("%s.SessionMaxTimeout", name), 0)
	cfg.ResponseMaxTimeout = config.GetDuration(fmt.Sprintf("%s.ResponseMaxTimeout", name), 0)
	cfg.RequestMaxTimeout = config.GetDuration(fmt.Sprintf("%s.RequestMaxTimeout", name), 0)
	cfg.SlowTimeout = config.GetDuration(fmt.Sprintf("%s.SlowTimeout", name), 0)
	cfg.CountTime = config.GetBool(fmt.Sprintf("%s.CountTime", name), false)
	cfg.RedialTimes = config.GetInt(fmt.Sprintf("%s.RedialTimes", name), 0)
	cfg.RedialInterval = config.GetDuration(fmt.Sprintf("%s.RedialTimes", name), 0)

	debug := config.GetBool("Debug", false)
	cfg.PrintDetail = debug
	// 命令行参数覆盖配置文件
	if len(parsers) > 0 {
		parser := parsers[0]
		cfg.ListenAddress = fmt.Sprintf("%s:%d",
			parser.GetOptVar("host", host).String(),
			parser.GetOptVar("port", port).Int(),
		)
		cfg.Network = parser.GetOptVar(fmt.Sprintf("%s.Network", name), cfg.Network).String()
	}

	return cfg
}

// DefaultBoxConf 创建运行配置
func DefaultBoxConf(parser *gcmd.Parser, config *gcfg.Config) *BoxConf {
	cfg := NewBoxConf("default.sandbox", config, parser)
	return cfg
}

// ListenPort 获取监听端口
func (that *BoxConf) ListenPort() string {
	_, port, err := net.SplitHostPort(that.ListenAddress)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	return port
}

// InnerIpPort 获取内网的服务地址
func (that *BoxConf) InnerIpPort() string {
	host, err := gipv4.GetIntranetIp()
	if err != nil {
		logger.Fatalf("%v", err)
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
func (that *BoxConf) EndpointConfig() drpc.EndpointConfig {
	var c = drpc.EndpointConfig{
		PrintDetail:       that.PrintDetail,
		Network:           that.Network,
		DefaultSessionAge: that.SessionMaxTimeout,
		DefaultContextAge: that.ResponseMaxTimeout,
		SlowCometDuration: that.SlowTimeout,
		CountTime:         that.CountTime,
	}
	if len(that.ListenAddress) > 0 {
		host, port, err := net.SplitHostPort(that.ListenAddress)
		if err != nil {
			logger.Fatalf("%v", err)
		}
		listenPort, _ := strconv.Atoi(port)
		c.LocalIP = host
		c.ListenPort = uint16(listenPort)
	}
	return c
}
