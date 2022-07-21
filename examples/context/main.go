package main

import (
	"github.com/osgochina/dmicro/drpc"
)

type TestCall struct {
	drpc.CallCtx
}

func (that *TestCall) Test(arg *string) (string, *drpc.Status) {
	that.Input()  // 当前请求的消息
	that.Output() // 当前请求的响应消息
	that.Swap()   // 当前session的交换空间
	that.Seq()
	return *arg, nil
}

type TestPush struct {
	drpc.PushCtx
}

func main() {
	endpoint := drpc.NewEndpoint(drpc.EndpointConfig{})
	endpoint.RouteCall(new(TestCall))
}
