package main

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
)

// DefaultSandBox  默认的服务
type DefaultSandBox struct {
	dserver.BaseSandbox
}

func (that *DefaultSandBox) Name() string {
	return "DefaultSandBox"
}

func (that *DefaultSandBox) Setup() error {
	fmt.Println("DefaultSandBox Setup")
	return nil
}

func (that *DefaultSandBox) Shutdown() error {
	fmt.Println("DefaultSandBox Shutdown")
	return nil
}

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_foo")
	dserver.CloseCtl()
	dserver.Setup(func(svr *dserver.DServer) {

		err := svr.AddSandBox(new(DefaultSandBox), svr.NewService("rpc"))
		if err != nil {
			logger.Fatal(context.TODO(), err)
		}
	})
}
