package registry

import (
	"github.com/gogf/gf/container/garray"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/utils/backoff"
	"net"
	"sync"
	"time"
)

type PluginRegistry struct {
	registry   Registry
	service    *Service
	allApis    *garray.StrArray
	addr       net.Addr
	registered bool
	mu         sync.RWMutex
}

var _ drpc.AfterRegRouterPlugin = new(PluginRegistry)
var _ drpc.AfterListenPlugin = new(PluginRegistry)
var _ drpc.BeforeCloseEndpointPlugin = new(PluginRegistry)

// NewRegistryPlugin 创建服务注册插件
func NewRegistryPlugin(registry Registry, service *Service) *PluginRegistry {
	return &PluginRegistry{
		registry: registry,
		service:  service,
		allApis:  garray.NewStrArray(),
	}
}

func (that *PluginRegistry) Name() string {
	return "registry"
}

func (that *PluginRegistry) AfterRegRouter(handler *drpc.Handler) error {
	that.allApis.Append(handler.Name())
	return nil
}

// AfterListen 服务启动成功，监听成功后，进行服务注册
func (that *PluginRegistry) AfterListen(addr net.Addr) (err error) {
	that.mu.RLock()
	config := that.registry.Options()
	that.mu.RUnlock()
	regFunc := func(service *Service) error {
		ttl := DefaultRegisterTTL
		if config.Context != nil {
			ttl1, ok := config.Context.Value("RegisterTTL").(time.Duration)
			if ok {
				ttl = ttl1
			}
		}

		rOpts := []RegisterOption{RegisterTTL(ttl)}
		var regErr error

		for i := 0; i < 3; i++ {
			// 循环注册三次，直到注册成功，如果三次都为成功，则报错
			if e := that.registry.Register(service, rOpts...); e != nil {
				regErr = e
				time.Sleep(backoff.Do(i + 1))
				continue
			}
			// 如果执行三次，最后成功了，则清除错误
			regErr = nil
			break
		}

		return regErr
	}
	addr.String()
	node := &Node{
		Id:       that.service.Name + "-" + addr.String(),
		Address:  addr.String(),
		Metadata: make(map[string]string),
		Paths:    that.allApis.Slice(),
	}
	node.Metadata["registry"] = that.registry.String()
	node.Metadata["server"] = that.service.Name

	service := &Service{
		Name:    that.service.Name,
		Version: that.service.Version,
		Nodes:   []*Node{node},
	}
	that.addr = addr
	return regFunc(service)
}

// BeforeCloseEndpoint 关闭服务之前先取消注册
func (that *PluginRegistry) BeforeCloseEndpoint(endpoint drpc.Endpoint) error {
	node := &Node{
		Id:      that.service.Name + "-" + that.addr.String(),
		Address: that.addr.String(),
	}

	service := &Service{
		Name:    that.service.Name,
		Version: that.service.Version,
		Nodes:   []*Node{node},
	}
	err := that.registry.Deregister(service)
	that.registry.Stop()
	return err
}
