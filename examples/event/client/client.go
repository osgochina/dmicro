package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/plugin/event"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
	"time"
)

func main() {
	bus := eventbus.New(grand.S(8))
	cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true, RedialTimes: -1, RedialInterval: time.Second}, event.NewEventPlugin(bus))
	defer cli.Close()

	cli.RoutePush(new(Push))
	_ = bus.Listen(event.OnReceiveEvent, eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		onReceive := e.(*event.OnReceive)
		fmt.Println(onReceive.ReadCtx.Input().String())
		return nil
	}))
	_ = bus.Listen(event.OnConnectEvent, eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		onConnect := e.(*event.OnConnect)
		fmt.Println(onConnect.Session.RemoteAddr())
		return nil
	}))
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
	logger.Printf(context.TODO(), "Wait 10 seconds to receive the push...")
	time.Sleep(time.Second * 10)
}

type Push struct {
	drpc.PushCtx
}

func (that *Push) Status(arg *string) *drpc.Status {
	logger.Printf(context.TODO(), "%s", *arg)
	return nil
}
