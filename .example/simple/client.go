package main

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true})
	defer cli.Close()

	cli.RoutePush(new(Push))

	sess, stat := cli.Dial("127.0.0.1:9091")
	if !stat.OK() {
		glog.Fatalf("%v", stat)
	}
	var result int
	stat = sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
		message.WithSetMeta("author", "henrylee2cn"),
	).Status()
	if !stat.OK() {
		glog.Fatalf("%v", stat)
	}
	glog.Printf("result: %d", result)
	glog.Printf("Wait 10 seconds to receive the push...")
	time.Sleep(time.Second * 10)
}

type Push struct {
	drpc.PushCtx
}

func (that *Push) Status(arg *string) *drpc.Status {
	glog.Printf("%s", *arg)
	return nil
}
