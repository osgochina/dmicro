package main

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/registry/etcd"
	"time"
)

func main() {

	reg := etcd.NewRegistry(registry.AddrList("127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"),
		etcd.LeasesInterval(10*time.Second),
		etcd.RegisterTTL(20*time.Second),
	)

	rSvr := &registry.Service{Name: "testregistry", Version: "0.1"}
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		CountTime:   true,
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	}, ignorecase.NewIgnoreCase(),
		registry.NewRegistryPlugin(reg, rSvr),
	)

	svr.RouteCall(new(Math))
	//go func() {
	//	time.Sleep(time.Second * 10)
	//	svr.Close()
	//}()
	go func() {
		rSvr2 := &registry.Service{Name: "testregistry", Version: "0.1"}
		svr2 := drpc.NewEndpoint(drpc.EndpointConfig{
			CountTime:   true,
			LocalIP:     "127.0.0.1",
			ListenPort:  9092,
			PrintDetail: true,
		}, ignorecase.NewIgnoreCase(),
			registry.NewRegistryPlugin(reg, rSvr2),
		)

		svr2.RouteCall(new(Math))
		_ = svr2.ListenAndServe()
	}()
	_ = svr.ListenAndServe()
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
