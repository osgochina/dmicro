package main

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/status"
)

type TestCall struct {
	drpc.CallCtx
}

// 注册方法名为"/test_call/test"
func (that *TestCall) Test(arg *string) (string, *drpc.Status) {
	that.Input()  // 当前请求的消息
	that.Output() // 当前请求的响应消息
	that.Swap()   // 当前session的交换空间
	return *arg, nil
}

type TestPush struct {
	drpc.PushCtx
}

// 注册方法名为"/test_push/test"
func (that *TestPush) Test(arg *string) *drpc.Status {
	that.Swap() // 当前session的交换空间
	return nil
}

func main() {
	endpoint := drpc.NewEndpoint(drpc.EndpointConfig{})
	endpoint.RouteCall(new(TestCall))
	endpoint.RoutePush(new(TestPush))
	endpoint.SetUnknownPush(func(ctx drpc.UnknownPushCtx) *status.Status {
		fmt.Println("UnknownPush")
		ctx.Context()
		ctx.Swap()
		return nil
	})
	endpoint.SetUnknownCall(func(ctx drpc.UnknownCallCtx) (interface{}, *status.Status) {
		fmt.Println("UnknownCall")
		ctx.Context()
		ctx.Swap()
		return nil, nil
	})
}
