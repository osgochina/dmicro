package main

import (
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/examples/dserver/sandbox"
	"github.com/osgochina/dmicro/logger"
)

func main() {
	dserver.Setup(func(svr *dserver.DServer) {
		//svr.ProcessModel(dserver.ProcessModelMulti)
		svr.SetInheritListener([]dserver.InheritAddr{
			{Network: "tcp", Host: "127.0.0.1", Port: 8199, ServerName: "default"},
			{Network: "http", Host: "127.0.0.1", Port: 8080, ServerName: "ghttp1"},
		})
		err := svr.AddSandBox(new(sandbox.DefaultSandBox))
		if err != nil {
			logger.Fatal(err)
		}
		err = svr.AddSandBox(new(sandbox.DefaultSandBox1))
		if err != nil {
			logger.Fatal(err)
		}
	})
}
