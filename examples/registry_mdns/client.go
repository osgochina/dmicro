package main

import (
	"github.com/osgochina/dmicro/client"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {
	serviceName := "testregistry"
	//err := registry.DefaultRegistry.Init(registry.ServiceName(serviceName))
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//cli := client.NewRpcClient(serviceName, client.OptRegistry(registry.DefaultRegistry))
	cli := client.NewRpcClient(serviceName)
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
