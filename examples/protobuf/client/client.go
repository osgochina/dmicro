package main

import (
	"context"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{Network: "unix", PrintDetail: true, RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()

	cli.RoutePush(new(Push))

	sess, stat := cli.Dial("127.0.0.1:9091", pbproto.NewPbProtoFunc())
	if !stat.OK() {
		logger.Fatalf(context.TODO(), "%v", stat)
	}
	for i := 0; i < 100; i++ {
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
		logger.Printf(context.TODO(), "Wait 10 seconds to receive the push...")
		time.Sleep(time.Second * 1)
	}
}

type Push struct {
	drpc.PushCtx
}

func (that *Push) Status(arg *string) *drpc.Status {
	logger.Printf(context.TODO(), "%s", *arg)
	return nil
}
