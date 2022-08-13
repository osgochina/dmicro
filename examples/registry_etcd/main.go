package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/registry/etcd"
	"github.com/osgochina/dmicro/server"
	"time"
)

func main() {
	serviceName := "testregistry"
	serviceVersion := "1.0.1"

	svr := getOrigServer(serviceName, serviceVersion)
	go func() {
		time.Sleep(10 * time.Second)
		_ = svr.Close()
		reg1 := etcd.NewRegistry(registry.AddrList("127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"),
			registry.LeasesInterval(10*time.Second),
			etcd.RegisterTTL(20*time.Second),
			registry.ServiceName("testregistry"),
			registry.ServiceVersion("0.2"),
		)
		svr2 := drpc.NewEndpoint(drpc.EndpointConfig{
			CountTime:   true,
			LocalIP:     "127.0.0.1",
			ListenPort:  9092,
			PrintDetail: true,
		}, ignorecase.NewIgnoreCase(),
			registry.NewRegistryPlugin(reg1),
		)

		svr2.RouteCall(new(Math))
		_ = svr2.ListenAndServe()
	}()
	_ = svr.ListenAndServe()
	select {}
}

func getOrigServer(serviceName, serviceVersion string) drpc.Endpoint {
	etcd.SetPrefix("/vprix/registry/dev/")
	reg := etcd.NewRegistry(
		registry.ServiceName(serviceName),
		registry.ServiceVersion(serviceVersion),
		registry.AddrList("127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"),
		registry.LeasesInterval(10*time.Second),
		etcd.RegisterTTL(20*time.Second),
	)
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		CountTime:   true,
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	}, ignorecase.NewIgnoreCase(),
		registry.NewRegistryPlugin(reg),
	)

	svr.RouteCall(new(Math))
	return svr
}

func getRpcServer(serviceName, serviceVersion string) *server.RpcServer {
	reg := etcd.NewRegistry(
		registry.ServiceName(serviceName),
		registry.ServiceVersion(serviceVersion),
		registry.AddrList("127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"),
		registry.LeasesInterval(10*time.Second),
		etcd.RegisterTTL(20*time.Second),
	)
	svr := server.NewRpcServer(serviceName,
		server.OptServiceVersion(serviceVersion),
		server.OptRegistry(reg),
		server.OptListenAddress("127.0.0.1:9091"),
		server.OptCountTime(true),
		server.OptPrintDetail(true),
	)
	svr.RouteCall(new(Math))
	_ = svr.ListenAndServe()
	return svr
}

type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// test meta
	logger.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
