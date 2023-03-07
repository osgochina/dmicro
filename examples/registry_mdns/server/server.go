package main

import (
	"context"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/server"
)

func main() {
	serviceName := "testregistry"
	serviceVersion := "1.0.0"
	reg := registry.DefaultRegistry
	err := reg.Init(registry.OptServiceName(serviceName), registry.OptServiceVersion(serviceVersion))
	if err != nil {
		logger.Fatal(context.TODO(), err)
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
	logger.Infof(context.TODO(), "author: %s", that.PeekMeta("author"))
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}
