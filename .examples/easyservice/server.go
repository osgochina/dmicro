package main

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/osgochina/dmicro/.examples/easyservice/sandbox"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/easyservice"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils"
	"github.com/osgochina/dmicro/utils/graceful"
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
		tlsCertFile := fmt.Sprintf("%s/../quic/cert.pem", gfile.MainPkgPath())
		tlsKeyFile := fmt.Sprintf("%s/../quic/key.pem", gfile.MainPkgPath())
		tlsConfig, err := utils.NewTLSConfigFromFile(tlsCertFile, tlsKeyFile)

		// 使用master worker 进程模型实现平滑重启
		err = graceful.SetInheritListener([]graceful.InheritAddr{
			{Network: "tcp", Host: "127.0.0.1", Port: "8199"},
			{Network: "quic", Host: "127.0.0.1", Port: "8198", TlsConfig: tlsConfig},
			{Network: "http", Host: "127.0.0.1", Port: "8080", ServerName: "default"},
		})
		if err != nil {
			logger.Error(err)
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
		rpc2.Endpoint().SetTLSConfig(tlsConfig)
		svr.AddSandBox(rpc2)
		http := sandbox.NewHttpSandBox(svr)
		svr.AddSandBox(http)
	})
}
