package sandbox

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/metrics"
	"github.com/osgochina/dmicro/metrics/prometheus"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/server"
)

type AdminSandBox struct {
	dserver.ServiceSandbox
	rpcServer *server.RpcServer
}

func Login(ctx drpc.CallCtx, args *struct{}) (string, *drpc.Status) {
	return "login", nil
}
func (that *AdminSandBox) Name() string {
	return "AdminSandBox"
}

func (that *AdminSandBox) Setup() error {
	serviceName := "user"
	serviceVersion := "1.0.0"
	opts := that.Config.RpcServerOption(that.Name())
	opts = append(opts,
		server.OptEnableHeartbeat(true),
		server.OptListenAddress("127.0.0.1:8199"),
		server.OptRegistry(
			registry.NewRegistry(
				registry.OptServiceName(serviceName),
				registry.OptServiceVersion(serviceVersion),
			),
		),
		server.OptMetrics(prometheus.NewPromMetrics(
			metrics.OptHost("0.0.0.0"),
			metrics.OptPort(9101),
			metrics.OptPath("/metrics"),
			metrics.OptServiceName("test_one"),
		)),
	)
	that.rpcServer = server.NewRpcServer(that.Name(), opts...)
	that.rpcServer.SubRoute("/app").RouteCallFunc(Login)
	return that.rpcServer.ListenAndServe()
}

func (that *AdminSandBox) Shutdown() error {
	fmt.Println("AdminSandBox Shutdown")
	that.rpcServer.Close()
	return nil
}
