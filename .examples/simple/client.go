package main

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true, RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()

	cli.RoutePush(new(Push))

	sess, stat := cli.Dial("127.0.0.1:9091")
	if !stat.OK() {
		logger.Fatalf("%v", stat)
	}
	for i := 0; i < 100; i++ {
		var result int
		stat = sess.Call("/math/add",
			[]int{1, 2, 3, 4, 5},
			&result,
			message.WithSetMeta("author", "liuzhiming"),
		).Status()
		if !stat.OK() {
			logger.Fatalf("%v", stat)
		}
		logger.Printf("result: %d", result)
		logger.Printf("Wait 10 seconds to receive the push...")
		time.Sleep(time.Second * 1)
	}
}

type Push struct {
	drpc.PushCtx
}

func (that *Push) Status(arg *string) *drpc.Status {
	logger.Printf("%s", *arg)
	return nil
}
