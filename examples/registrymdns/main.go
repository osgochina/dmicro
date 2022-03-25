package main

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
)

func main() {
	err := registry.DefaultRegistry.Init(
		registry.ServiceName("testregistry"),
		registry.ServiceVersion("0.1"),
	)
	if err != nil {
		logger.Fatal(err)
	}
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		CountTime:   true,
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	}, ignorecase.NewIgnoreCase(),
		registry.NewRegistryPlugin(registry.DefaultRegistry),
	)

	svr.RouteCall(new(Math))
	_ = svr.ListenAndServe()
	select {}
}

type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// test meta
	glog.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
