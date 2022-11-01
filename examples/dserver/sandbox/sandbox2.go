package sandbox

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/osgochina/dmicro/dserver"
)

type DefaultSandBox1 struct {
	dserver.BaseSandbox
	svr *ghttp.Server
}

func (that *DefaultSandBox1) Name() string {
	return "DefaultSandBox1"
}
func (that *DefaultSandBox1) Abc() string {
	return "DefaultSandBox1"
}
func (that *DefaultSandBox1) Setup() error {
	fmt.Println("DefaultSandBox1 Setup")
	c := g.Cfg()
	fmt.Println(c.Get("sandbox"))
	//that.svr = g.Server("ghttp1")
	//that.svr.BindHandler("/", func(r *ghttp.Request) {
	//	time.Sleep(10 * time.Second)
	//	r.Response.WriteExit("hello world!", "pid:"+gconv.String(os.Getpid()))
	//})
	//that.svr.SetPort(8080)
	//return that.svr.Start()
	return nil
}

func (that *DefaultSandBox1) Shutdown() error {
	fmt.Println("DefaultSandBox1 Shutdown")

	return nil
	//return that.svr.Shutdown()
}
