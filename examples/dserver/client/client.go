package main

import (
	"context"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {
	ctx := context.Background()

	cli := drpc.NewEndpoint(drpc.EndpointConfig{Network: "tcp", RedialTimes: 1})
	defer cli.Close()

	sess, stat := cli.Dial("127.0.0.1:8199")
	if !stat.OK() {
		logger.Fatalf(ctx, "%v", stat)
	}

	for true {
		var result string
		stat = sess.Call("/app/home",
			nil,
			&result,
		).Status()
		if !stat.OK() {
			logger.Fatalf(ctx, "%v", stat)
		}
		logger.Printf(ctx, "result: %s", result)
		time.Sleep(time.Second)
	}
}
