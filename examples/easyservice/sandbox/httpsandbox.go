package sandbox

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/osgochina/dmicro/easyservice"
)

// HttpSandBox  默认的服务
type HttpSandBox struct {
	id      int
	name    string
	service *easyservice.EasyService
	httpSvr *ghttp.Server
}

// NewHttpSandBox 创建一个默认的服务沙盒
func NewHttpSandBox(service *easyservice.EasyService) *HttpSandBox {
	id := easyservice.GetNextSandBoxId()

	sBox := &HttpSandBox{
		id:      id,
		service: service,
	}
	sBox.httpSvr = ghttp.GetServer("default")
	sBox.httpSvr.SetPort(8080)
	sBox.httpSvr.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("ok")
	})
	return sBox
}

func (that *HttpSandBox) ID() int {
	return that.id
}

func (that *HttpSandBox) Name() string {
	return that.name
}

func (that *HttpSandBox) Setup() error {
	return that.httpSvr.Start()
}

func (that *HttpSandBox) Shutdown() error {
	return that.httpSvr.Shutdown()
}

func (that *HttpSandBox) Service() *easyservice.EasyService {
	return that.service
}
