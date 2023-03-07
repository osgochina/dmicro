package main

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true, RedialTimes: 1, RedialInterval: time.Second})
	defer cli.Close()

	sess, stat := cli.Dial("127.0.0.1:9091")
	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}
	n := 1
	for {
		var result string
		stat = sess.Call("/grace/sleep",
			5,
			&result,
		).Status()
		if !stat.OK() {
			logger.Error(context.TODO(), stat.Cause())
			time.Sleep(1 * time.Second)
		}
		fmt.Printf("%d.%s\n", n, result)
		n++
		time.Sleep(1 * time.Second)
	}

}
