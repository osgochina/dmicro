package main

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/graceful"
	"time"
)

func main() {

	err := graceful.SetInheritListener([]graceful.InheritAddr{
		{Network: "tcp", Host: "127.0.0.1", Port: "9091"},
		{Network: "tcp", Host: "127.0.0.1", Port: "9092"},
		{Network: "tcp", Host: "127.0.0.1", Port: "9093"},
		{Network: "http", Host: "127.0.0.1", Port: "8080", ServerName: "test"},
		{Network: "http", Host: "127.0.0.1", Port: "8081", ServerName: "test"},
		{Network: "http", Host: "127.0.0.1", Port: "8082", ServerName: "test"},
	})
	if err != nil {
		logger.Error(context.TODO(), err)
		return
	}
	go graceful.GraceSignal()
	gSvr := g.Server("test")
	gSvr.SetPort(8080, 8081, 8082)
	gSvr.BindHandler("/", func(r *ghttp.Request) {
		str := "start:" + gtime.Now().String()
		time.Sleep(10 * time.Second)
		str += "\nend:" + gtime.Now().String()
		r.Response.Write("hello: \n" + str)
	})
	gSvr.Shutdown()
	gSvr.Start()

	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	})

	svr.RouteCall(new(Grace))
	go svr.ListenAndServe()

	svr1 := drpc.NewEndpoint(drpc.EndpointConfig{
		LocalIP:     "127.0.0.1",
		ListenPort:  9092,
		PrintDetail: true,
	})

	svr1.RouteCall(new(Grace))
	go svr1.ListenAndServe()
	svr2 := drpc.NewEndpoint(drpc.EndpointConfig{
		LocalIP:     "127.0.0.1",
		ListenPort:  9093,
		PrintDetail: true,
	})

	svr2.RouteCall(new(Grace))
	logger.Warning(context.TODO(), svr2.ListenAndServe())
}

type Grace struct {
	drpc.CallCtx
}

func (m *Grace) Sleep(arg *int) (string, *drpc.Status) {
	logger.Infof(context.TODO(), "sleep %d", *arg)
	if *arg > 0 {
		sleep := *arg
		time.Sleep(time.Duration(sleep) * time.Second)
	}
	// response
	return "sleep", nil
}
