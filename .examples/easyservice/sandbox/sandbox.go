package sandbox

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/easyservice"
)

// DefaultSandBox  默认的服务
type DefaultSandBox struct {
	id       int
	name     string
	boxConf  *easyservice.BoxConf
	service  *easyservice.EasyService
	endpoint drpc.Endpoint
}

// NewDefaultSandBox 创建一个默认的服务沙盒
func NewDefaultSandBox(cfg *easyservice.BoxConf, globalLeftPlugin ...drpc.Plugin) *DefaultSandBox {
	id := easyservice.GetNextSandBoxId()
	if len(cfg.SandBoxName) <= 0 {
		cfg.SandBoxName = fmt.Sprintf("default_%d", id)
	}

	sBox := &DefaultSandBox{
		id:      id,
		name:    cfg.SandBoxName,
		boxConf: cfg,
	}
	var pluginArray []drpc.Plugin
	if len(globalLeftPlugin) > 0 {
		pluginArray = append(pluginArray, globalLeftPlugin...)
	}
	sBox.endpoint = drpc.NewEndpoint(cfg.EndpointConfig(), pluginArray...)
	return sBox
}

func (that *DefaultSandBox) ID() int {
	return that.id
}

func (that *DefaultSandBox) Name() string {
	return that.name
}

func (that *DefaultSandBox) Setup() error {

	return that.endpoint.ListenAndServe()
}

func (that *DefaultSandBox) Shutdown() error {
	return that.endpoint.Close()
}

func (that *DefaultSandBox) Endpoint() drpc.Endpoint {
	return that.endpoint
}

func (that *DefaultSandBox) Service() *easyservice.EasyService {
	return that.service
}
