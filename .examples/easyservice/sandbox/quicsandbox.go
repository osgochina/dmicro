package sandbox

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/easyservice"
)

// QUICSandBox  默认的服务
type QUICSandBox struct {
	id       int
	name     string
	boxConf  *easyservice.BoxConf
	service  *easyservice.EasyService
	endpoint drpc.Endpoint
}

// NewQUICSandBox 创建一个默认的服务沙盒
func NewQUICSandBox(cfg *easyservice.BoxConf, globalLeftPlugin ...drpc.Plugin) *DefaultSandBox {
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

func (that *QUICSandBox) ID() int {
	return that.id
}

func (that *QUICSandBox) Name() string {
	return that.name
}

func (that *QUICSandBox) Setup() error {

	return that.endpoint.ListenAndServe()
}

func (that *QUICSandBox) Shutdown() error {
	return that.endpoint.Close()
}

func (that *QUICSandBox) Endpoint() drpc.Endpoint {
	return that.endpoint
}

func (that *QUICSandBox) Service() *easyservice.EasyService {
	return that.service
}
