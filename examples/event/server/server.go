package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/event"
	"github.com/osgochina/dmicro/eventbus"
)

func main() {
	bus := eventbus.New(grand.S(8))
	endpoint := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenIP:    "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: false,
	}, event.NewEventPlugin(bus))
	endpoint.RouteCall(new(Math))
	_ = bus.Listen(event.OnAcceptEvent, eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		onAccept := e.(*event.OnAccept)
		fmt.Println(onAccept.Session.RemoteAddr())
		return nil
	}))
	_ = bus.Listen(event.OnReceiveEvent, eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		onReceive := e.(*event.OnReceive)
		fmt.Println(onReceive.ReadCtx.Input().String())
		return nil
	}))
	_ = bus.Listen(event.OnCloseEvent, eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		onClose := e.(*event.OnClose)
		fmt.Println(onClose.Session.RemoteAddr())
		return nil
	}))
	endpoint.ListenAndServe()
}

type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// test meta
	glog.Infof(context.TODO(), "author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
