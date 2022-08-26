package sandbox

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/server"
)

// DefaultSandBox  默认的服务
type DefaultSandBox struct {
	dserver.ServiceSandbox
	rpcServer *server.RpcServer
	TestName  string
}

func Home(ctx drpc.CallCtx, args *struct{}) (string, *drpc.Status) {
	return "home", nil
}
func (that *DefaultSandBox) Name() string {
	return "DefaultSandBox"
}

func (that *DefaultSandBox) Setup() error {
	that.rpcServer = server.NewRpcServer(that.Name(), that.Config.RpcServerOption(that.Name())...)
	that.rpcServer.SubRoute("/app").RouteCallFunc(Home)
	return that.rpcServer.ListenAndServe()
}

func (that *DefaultSandBox) Shutdown() error {
	fmt.Println("DefaultSandBox Shutdown")
	that.rpcServer.Close()
	return nil
}
