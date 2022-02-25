package main

import (
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/registry/etcd"
	"time"
)

func main() {

	reg := etcd.NewRegistry(registry.AddrList("127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"))

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

	// broadcast per 5s
	go func() {
		for {
			time.Sleep(time.Second * 5)
			s, e := reg.ListServices()
			if e != nil {
				logger.Error(e)
			}
			logger.Info(s)
			svr.RangeSession(func(sess drpc.Session) bool {
				sess.Push(
					"/push/status",
					fmt.Sprintf("this is a broadcast, server time: %v", time.Now()),
				)
				return true
			})
		}
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
