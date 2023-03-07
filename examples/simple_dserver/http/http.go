package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
	"os"
)

// HttpSandBox  默认的服务
type HttpSandBox struct {
	dserver.BaseSandbox
	http *ghttp.Server
}

func (that *HttpSandBox) Name() string {
	return "HttpSandBox"
}

func (that *HttpSandBox) Setup() error {
	fmt.Println("HttpSandBox Setup")
	that.http = g.Server("ghttp")
	that.http.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteExit("hello world!", "pid:"+gconv.String(os.Getpid()), "\n")
	})
	that.http.SetPort(8080)
	return that.http.Start()
}

func (that *HttpSandBox) Shutdown() error {
	fmt.Println("HttpSandBox Shutdown")
	return that.http.Shutdown()
}

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_http")
	dserver.Setup(func(svr *dserver.DServer) {
		err := svr.AddSandBox(new(HttpSandBox))
		if err != nil {
			logger.Fatal(context.TODO(), err)
		}
	})
}
