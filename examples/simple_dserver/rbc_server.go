package main

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
)

// DefaultDRpcSandBox  默认的服务
type DefaultDRpcSandBox struct {
	dserver.BaseSandbox
	endpoint drpc.Endpoint
}

func (that *DefaultDRpcSandBox) Name() string {
	return "DefaultDRpcSandBox"
}

func (that *DefaultDRpcSandBox) Setup() error {
	fmt.Println("DefaultDRpcSandBox Setup")
	cfg := that.Config.EndpointConfig(that.Name())
	cfg.ListenPort = 9091
	cfg.PrintDetail = true
	that.endpoint = drpc.NewEndpoint(cfg)
	that.endpoint.RouteCall(new(Math))
	return that.endpoint.ListenAndServe()
}

func (that *DefaultDRpcSandBox) Shutdown() error {
	fmt.Println("DefaultDRpcSandBox Shutdown")
	return that.endpoint.Close()
}

// Math rpc请求的最终处理器，必须集成drpc.CallCtx
type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// test meta
	logger.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_drpc")
	dserver.Setup(func(svr *dserver.DServer) {
		err := svr.AddSandBox(new(DefaultDRpcSandBox))
		if err != nil {
			logger.Fatal(err)
		}
	})
}
