package main

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/osgochina/dmicro/.examples/easyservice/sandbox"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/easyservice"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils"
)

func Home(ctx drpc.CallCtx, args *struct{}) (string, *drpc.Status) {
	return "home", nil
}

func main() {
	logger.SetDebug(true)
	easyservice.Setup(func(svr *easyservice.EasyService) {
		//注册服务停止时要执行法方法
		svr.BeforeStop(func(service *easyservice.EasyService) bool {
			fmt.Println("server stop")
			return true
		})
		tlsConfigQUIC, err := utils.NewTLSConfigFromFile(fmt.Sprintf("%s/../quic/cert.pem", gfile.MainPkgPath()), fmt.Sprintf("%s/../quic/key.pem", gfile.MainPkgPath()))
		tlsConfigKCP, err := utils.NewTLSConfigFromFile(fmt.Sprintf("%s/../kcp/cert.pem", gfile.MainPkgPath()), fmt.Sprintf("%s/../kcp/key.pem", gfile.MainPkgPath()))

		//// 使用master worker 进程模型实现平滑重启
		//err = graceful.SetInheritListener([]graceful.InheritAddr{
		//	{Network: svr.CmdParser().GetOptVar("network", "tcp").String(), Host: svr.CmdParser().GetOptVar("host", "127.0.0.1").String(), Port: "8199"},
		//	{Network: "http", Host: "127.0.0.1", Port: "8080", ServerName: "default"},
		//})
		if err != nil {
			logger.Fatal(err)
			return
		}
		var cfg = easyservice.DefaultBoxConf(svr.CmdParser(), svr.Config())
		rpc := sandbox.NewDefaultSandBox(cfg)
		rpc.Endpoint().SubRoute("/app").RouteCallFunc(Home)
		svr.AddSandBox(rpc)

		var cfg2 = easyservice.DefaultBoxConf(svr.CmdParser(), svr.Config())
		cfg2.ListenAddress = "127.0.0.1:8198"
		cfg2.Network = "quic"
		rpc2 := sandbox.NewQUICSandBox(cfg2)
		rpc2.Endpoint().SubRoute("/app").RouteCallFunc(Home)
		rpc2.Endpoint().SetTLSConfig(tlsConfigQUIC)
		svr.AddSandBox(rpc2)

		var cfg3 = easyservice.DefaultBoxConf(svr.CmdParser(), svr.Config())
		cfg3.ListenAddress = "127.0.0.1:8197"
		cfg3.Network = "kcp"
		rpc3 := sandbox.NewKCPSandBox(cfg3)
		rpc3.Endpoint().SubRoute("/app").RouteCallFunc(Home)
		rpc3.Endpoint().SetTLSConfig(tlsConfigKCP)
		svr.AddSandBox(rpc3)

		var cfg4 = easyservice.NewBoxConf("sandboxBak", svr.Config())
		rpc4 := sandbox.NewDefaultSandBox(cfg4)
		rpc4.Endpoint().SubRoute("/app").RouteCallFunc(Home)
		svr.AddSandBox(rpc4)

		http := sandbox.NewHttpSandBox(svr)
		svr.AddSandBox(http)
	})
}
