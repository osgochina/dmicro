package sandbox

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/easyservice"
)

// KCPSandBox  默认的服务
type KCPSandBox struct {
	id       int
	name     string
	boxConf  *easyservice.BoxConf
	service  *easyservice.EasyService
	endpoint drpc.Endpoint
}

// NewKCPSandBox 创建一个默认的服务沙盒
func NewKCPSandBox(cfg *easyservice.BoxConf, globalLeftPlugin ...drpc.Plugin) *KCPSandBox {
	id := easyservice.GetNextSandBoxId()
	if len(cfg.SandBoxName) <= 0 {
		cfg.SandBoxName = fmt.Sprintf("default_%d", id)
	}

	sBox := &KCPSandBox{
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

func (that *KCPSandBox) ID() int {
	return that.id
}

func (that *KCPSandBox) Name() string {
	return that.name
}

func (that *KCPSandBox) Setup() error {

	return that.endpoint.ListenAndServe()
}

func (that *KCPSandBox) Shutdown() error {
	return that.endpoint.Close()
}

func (that *KCPSandBox) Endpoint() drpc.Endpoint {
	return that.endpoint
}

func (that *KCPSandBox) Service() *easyservice.EasyService {
	return that.service
}
