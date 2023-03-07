package sandbox

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/osgochina/dmicro/dserver"
)

type UserSandBox struct {
	dserver.ServiceSandbox
	svr *ghttp.Server
}

func (that *UserSandBox) Name() string {
	return "UserSandBox"
}

func (that *UserSandBox) Setup() error {
	that.svr = g.Server()
	return that.svr.Start()
}

func (that *UserSandBox) Shutdown() error {
	fmt.Println("UserSandBox Shutdown")

	return nil
}
