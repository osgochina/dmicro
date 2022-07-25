package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"time"
)

func main() {
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	},
	)
	svr.PluginContainer().AppendRight(ignorecase.NewIgnoreCase()) //使用IgnoreCase插件，活跃serviceName的大小写
	go svr.ListenAndServe()

	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{PrintDetail: false},
	)
	cli.PluginContainer().AppendRight(ignorecase.NewIgnoreCase()) // 使用IgnoreCase插件，活跃serviceName的大小写
	cli.Dial(":9090")
	time.Sleep(time.Second * 20)
}
