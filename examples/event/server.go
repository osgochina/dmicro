package main

import (
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/event"
	"github.com/osgochina/dmicro/eventbus"
)

func main() {
	endpoint := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenIP:    "127.0.0.1",
		ListenPort:  9091,
		PrintDetail: false,
		CountTime:   false,
	}, event.NewEventPlugin())
	endpoint.RouteCall(new(Math))
	_ = endpoint.EventBus().Listen(event.OnAcceptEvent, eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		onAccept := e.(*event.OnAccept)
		fmt.Println(onAccept.Session.RemoteAddr())
		return nil
	}))
	_ = endpoint.EventBus().Listen(event.OnReceiveEvent, eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		onReceive := e.(*event.OnReceive)
		fmt.Println(onReceive.ReadCtx.Input().String())
		return nil
	}))
	_ = endpoint.EventBus().Listen(event.OnCloseEvent, eventbus.ListenerFunc(func(e eventbus.IEvent) error {
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
	glog.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
