package main

import (
	"github.com/osgochina/dmicro/client"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/registry/etcd"
	"github.com/osgochina/dmicro/selector"
	"time"
)

func main() {
	serviceName := "testregistry"
	etcd.SetPrefix("/vprix/registry/dev/")
	reg := etcd.NewRegistry(registry.AddrList("127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"),
		registry.ServiceName(serviceName),
	)
	sel := selector.NewSelector(selector.Registry(reg))
	cli := client.NewRpcClient(serviceName, client.OptSelector(sel))
	for {
		var result int
		stat := cli.Call("/math/add",
			[]int{1, 2, 3, 4, 5},
			&result,
			message.WithSetMeta("author", "liuzhiming"),
		).Status()
		if !stat.OK() {
			logger.Fatalf("%v", stat)
		}
		logger.Printf("result: %d", result)
		time.Sleep(time.Second * 2)
	}

}
