package main

import (
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/examples/dserver/sandbox"
	"github.com/osgochina/dmicro/logger"
)

func main() {
	dserver.Setup(func(svr *dserver.DServer) {
		svr.MultiProcess(true)
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
