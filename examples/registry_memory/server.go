package main

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/ignorecase"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/registry/memory"
	"time"
)

func main() {
	reg := memory.NewRegistry(
		registry.OptServiceName("testregistry"),
		registry.OptServiceVersion("1.0.0"),
	)
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		LocalIP:     "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: true,
	}, ignorecase.NewIgnoreCase(),
		registry.NewRegistryPlugin(reg),
	)
	svr.RouteCall(new(Math))
	go func() {
		time.Sleep(3 * time.Second)
		s, err := reg.GetService("testregistry")
		if err != nil {
			logger.Error(err)
		}
		for _, s1 := range s {
			for _, n := range s1.Nodes {
				fmt.Printf("node id: %s\n", n.Id)
				fmt.Printf("node address: %s\n", n.Address)
			}

		}
	}()
	_ = svr.ListenAndServe()
	select {}
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
