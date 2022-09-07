package dserver

import (
	"fmt"
	"github.com/gogf/gf/os/gcfg"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/server"
	"time"
)

// Config 配置文件对象
type Config struct {
	*gcfg.Config
}

// EndpointConfig 通过配置文件获取配置信息
func (that *Config) EndpointConfig(sandboxName ...string) drpc.EndpointConfig {
	cfg := drpc.EndpointConfig{
		Network:           "tcp",                          //网络协议
		DefaultBodyCodec:  drpc.DefaultBodyCodec().Name(), // 默认的消息体编码格式
		DefaultSessionAge: 0,                              // 永久有效
		DefaultContextAge: 0,                              // 永久有效
		SlowCometDuration: 0,                              // 不记录
		PrintDetail:       false,                          // 是否打印日志详情

		ListenIP:   "0.0.0.0", //作为服务端角色时，要监听的服务器本地IP
		ListenPort: 0,         //作为服务端角色时，需要监听的本地端口号

		LocalIP:        "0.0.0.0",              //作为客户端角色时,请求服务端时候，本地使用的地址
		LocalPort:      0,                      //作为客户端角色时,请求服务端时候，本地使用的地址端口号
		DialTimeout:    30 * time.Second,       // 作为客户端角色时，请求服务端的超时时间
		RedialTimes:    1,                      // 仅限客户端角色使用,链接中断时候，试图链接服务端的最大重试次数。
		RedialInterval: time.Millisecond * 100, //仅限客户端角色使用 试图链接服务端时候，重试的时间间隔.

	}
	if len(sandboxName) > 0 {
		name := sandboxName[0]
		cj := that.Config.GetJson(fmt.Sprintf("sandbox.%s", name))
		if cj.IsNil() {
			return cfg
		}
		cfg.Network = cj.GetString("Network", "tcp")
		cfg.ListenIP = cj.GetString("ListenIP", "0.0.0.0")
		cfg.ListenPort = cj.GetUint16("ListenPort", 0)
		cfg.LocalIP = cj.GetString("LocalIP", "0.0.0.0")
		cfg.LocalPort = cj.GetUint16("LocalPort", 0)

		cfg.DefaultBodyCodec = cj.GetString("DefaultBodyCodec", drpc.DefaultBodyCodec().Name())
		cfg.DefaultSessionAge = time.Duration(cj.GetInt("DefaultSessionAge", 0)) * time.Second
		cfg.DefaultContextAge = time.Duration(cj.GetInt("DefaultContextAge", 0)) * time.Second
		cfg.SlowCometDuration = time.Duration(cj.GetInt("SlowCometDuration", 0)) * time.Second
		cfg.PrintDetail = cj.GetBool("PrintDetail", false)

		cfg.DialTimeout = time.Duration(cj.GetInt("DialTimeout", 0)) * time.Second
		cfg.RedialInterval = time.Duration(cj.GetInt("RedialInterval", 0)) * time.Second
		cfg.RedialTimes = cj.GetInt("RedialTimes", 1)
	}

	return cfg
}

// RpcServerOption 获取rpc server的参数
func (that *Config) RpcServerOption(serverName string) []server.Option {

	var opts []server.Option
	cfg := that.Config.GetJson(serverName)
	if cfg.IsNil() {
		return opts
	}

	opts = append(opts, server.OptNetwork(cfg.GetString("Network", "tcp")))
	opts = append(opts, server.OptListenAddress(
		fmt.Sprintf("%s:%d", cfg.GetString("ListenIP", "0.0.0.0"), cfg.GetUint16("ListenPort", 0)),
	))
	opts = append(opts, server.OptBodyCodec(cfg.GetString("DefaultBodyCodec", drpc.DefaultBodyCodec().Name())))
	opts = append(opts, server.OptSessionAge(time.Duration(cfg.GetInt("DefaultSessionAge", 0))*time.Second))
	opts = append(opts, server.OptContextAge(time.Duration(cfg.GetInt("DefaultContextAge", 0))*time.Second))
	opts = append(opts, server.OptSlowCometDuration(time.Duration(cfg.GetInt("SlowCometDuration", 0))*time.Second))
	opts = append(opts, server.OptPrintDetail(cfg.GetBool("PrintDetail", false)))

	return opts
}
