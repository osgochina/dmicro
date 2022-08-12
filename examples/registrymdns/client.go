package main

import (
	"github.com/osgochina/dmicro/client"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/selector"
	"time"
)

func main() {
	serviceName := "testregistry"
	err := registry.DefaultRegistry.Init(registry.ServiceName(serviceName))
	if err != nil {
		logger.Fatal(err)
	}
	sel := selector.NewSelector(selector.Registry(registry.DefaultRegistry))
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
