package main

import (
	"context"
	"github.com/osgochina/dmicro/client"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"time"
)

func main() {
	svr := &registry.Service{
		Nodes: []*registry.Node{
			{
				Address: "127.0.0.1:9091",
			},
		},
	}
	cli := client.NewRpcClient("testregistry", client.OptCustomService(svr))
	for {
		var result int
		stat := cli.Call("/math/add",
			[]int{1, 2, 3, 4, 5},
			&result,
			message.WithSetMeta("author", "liuzhiming"),
		).Status()
		if !stat.OK() {
			logger.Fatalf(context.TODO(), "%v", stat)
		}
		logger.Printf(context.TODO(), "result: %d", result)
		time.Sleep(time.Second * 2)
	}

}
