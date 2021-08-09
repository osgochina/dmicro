package easyserver

import (
	"fmt"
	"github.com/gogf/gf/net/gipv4"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"net"
	"strconv"
	"time"
)

type BoxConf struct {
	id                 int           //服务沙盒的id"
	SandBoxName        string        `json:"sandbox_name"   comment:"服务沙盒的名字"`
	Network            string        `json:"network"        comment:"使用的网络协议; tcp, tcp4, tcp6, unix or unixpacket"`
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
