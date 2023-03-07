package main

import (
	"context"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
)

//go:generate go build $GOFILE

func main() {

	srv := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort: 9090,
	})
	srv.RouteCall(new(math))
	srv.RoutePush(new(chat))
	srv.ListenAndServe()
}

type math struct {
	drpc.CallCtx
}

func (m *math) Add(arg *[]int) (int, *drpc.Status) {
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}

type chat struct {
	drpc.PushCtx
}

func (c *chat) Say(arg *string) *drpc.Status {
	logger.Printf(context.TODO(), "%s say: %q", c.PeekMeta("X-ID"), *arg)
	return nil
}
