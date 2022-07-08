package sandbox

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/dserver"
)

// DefaultSandBox  默认的服务
type DefaultSandBox struct {
	dserver.BaseSandbox
	endpoint drpc.Endpoint
	TestName string
}

func Home(ctx drpc.CallCtx, args *struct{}) (string, *drpc.Status) {
	return "home", nil
}
func (that *DefaultSandBox) Name() string {
	return "DefaultSandBox"
}

func (that *DefaultSandBox) Setup() error {
	fmt.Println("DefaultSandBox Setup")

	var c = drpc.EndpointConfig{
		PrintDetail: true,
		Network:     "tcp",
		LocalIP:     "127.0.0.1",
		ListenPort:  8199,
	}
	that.endpoint = drpc.NewEndpoint(c)
	that.endpoint.SubRoute("/app").RouteCallFunc(Home)
	return that.endpoint.ListenAndServe()
}

func (that *DefaultSandBox) Shutdown() error {
	fmt.Println("DefaultSandBox Shutdown")
	return that.endpoint.Close()
}
