package dserver

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gcfg"
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
		cj := that.Config.MustGet(context.TODO(), fmt.Sprintf("sandbox.%s", name))
		if cj.IsNil() {
			return cfg
		}
		cjJson := gjson.New(cj)
		cfg.Network = cjJson.Get("Network", "tcp").String()
		cfg.ListenIP = cjJson.Get("ListenIP", "0.0.0.0").String()
		cfg.ListenPort = cjJson.Get("ListenPort", 0).Uint16()
		cfg.LocalIP = cjJson.Get("LocalIP", "0.0.0.0").String()
		cfg.LocalPort = cjJson.Get("LocalPort", 0).Uint16()

		cfg.DefaultBodyCodec = cjJson.Get("DefaultBodyCodec", drpc.DefaultBodyCodec().Name()).String()
		cfg.DefaultSessionAge = time.Duration(cjJson.Get("DefaultSessionAge", 0).Int()) * time.Second
		cfg.DefaultContextAge = time.Duration(cjJson.Get("DefaultContextAge", 0).Int()) * time.Second
		cfg.SlowCometDuration = time.Duration(cjJson.Get("SlowCometDuration", 0).Int()) * time.Second
		cfg.PrintDetail = cjJson.Get("PrintDetail", false).Bool()

		cfg.DialTimeout = time.Duration(cjJson.Get("DialTimeout", 0).Int()) * time.Second
		cfg.RedialInterval = time.Duration(cjJson.Get("RedialInterval", 0).Int()) * time.Second
		cfg.RedialTimes = cjJson.Get("RedialTimes", 1).Int()
	}

	return cfg
}

// RpcServerOption 获取rpc server的参数
func (that *Config) RpcServerOption(serverName string) []server.Option {

	var opts []server.Option
	cfg := that.Config.MustGet(context.TODO(), serverName)
	if cfg == nil || cfg.IsNil() {
		return opts
	}
	cfgJson := gjson.New(cfg)
	opts = append(opts, server.OptNetwork(cfgJson.Get("Network", "tcp").String()))
	opts = append(opts, server.OptListenAddress(
		fmt.Sprintf("%s:%d", cfgJson.Get("ListenIP", "0.0.0.0").String(), cfgJson.Get("ListenPort", 0).Uint16()),
	))
	opts = append(opts, server.OptBodyCodec(cfgJson.Get("DefaultBodyCodec", drpc.DefaultBodyCodec().Name()).String()))
	opts = append(opts, server.OptSessionAge(time.Duration(cfgJson.Get("DefaultSessionAge", 0).Int())*time.Second))
	opts = append(opts, server.OptContextAge(time.Duration(cfgJson.Get("DefaultContextAge", 0).Int())*time.Second))
	opts = append(opts, server.OptSlowCometDuration(time.Duration(cfgJson.Get("SlowCometDuration", 0).Int())*time.Second))
	opts = append(opts, server.OptPrintDetail(cfgJson.Get("PrintDetail", false).Bool()))

	return opts
}
