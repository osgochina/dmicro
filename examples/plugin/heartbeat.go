package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/heartbeat"
	"time"
)

func main() {
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	},
		heartbeat.NewPong(), //使用pong插件，响应客户端的ping请求
	)
	go svr.ListenAndServe()

	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{PrintDetail: false},
		heartbeat.NewPing(3, true), // 使用ping插件，发送ping请求给服务端
	)
	cli.Dial(":9090")
	time.Sleep(time.Second * 20)
}
