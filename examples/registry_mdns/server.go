package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/server"
)

func main() {
	serviceName := "testregistry"
	serviceVersion := "1.0.0"
	reg := registry.DefaultRegistry
	err := reg.Init(registry.ServiceName(serviceName), registry.ServiceVersion(serviceVersion))
	if err != nil {
		logger.Fatal(err)
	}
	svr := server.NewRpcServer(serviceName,
		server.OptListenAddress("127.0.0.1:9091"),
		server.OptPrintDetail(true),
		server.OptServiceVersion(serviceVersion),
		server.OptRegistry(reg),
	)
	svr.RouteCall(new(Math))
	_ = svr.ListenAndServe()
}

type Math struct {
	drpc.CallCtx
}

func (that *Math) Add(arg *[]int) (int, *drpc.Status) {
	logger.Infof("author: %s", that.PeekMeta("author"))
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}
