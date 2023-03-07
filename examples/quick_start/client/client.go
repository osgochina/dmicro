package main

import (
	"context"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true, RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()

	sess, stat := cli.Dial("127.0.0.1:9091")
	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}
	var result int
	stat = sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
		message.WithSetMeta("author", "liuzhiming"),
	).Status()
	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}
	logger.Printf(context.TODO(), "result: %d", result)
}
