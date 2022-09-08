package registry

import (
	"errors"
	"github.com/gogf/gf/container/garray"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/backoff"
	"net"
	"sync"
	"time"
)

type PluginRegistry struct {
	registry       Registry
	service        *Service
	allApis        *garray.StrArray
	addr           net.Addr
	registered     bool
	mu             sync.RWMutex
	serviceVersion string
	serviceName    string
	exit           chan bool
}

var _ drpc.AfterRegRouterPlugin = new(PluginRegistry)
var _ drpc.AfterListenPlugin = new(PluginRegistry)
var _ drpc.BeforeCloseEndpointPlugin = new(PluginRegistry)

// NewRegistryPlugin 创建服务注册插件
func NewRegistryPlugin(registry Registry) *PluginRegistry {
	return &PluginRegistry{
		registry: registry,
		allApis:  garray.NewStrArray(),
		exit:     make(chan bool),
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
	that.addr = addr
	err = that.Register()

	go func() {
		t := time.NewTicker(that.leasesInterval())
		for {
			select {
			case <-t.C:
				that.mu.RLock()
				registered := that.registered
				that.mu.RUnlock()
				if registered {
					if err := that.Register(); err != nil {
						logger.Error(err)
					}
				}
			case <-that.exit:
				t.Stop()
				close(that.exit)
				return
			}
		}
	}()
	return err
}

// Register 注册
func (that *PluginRegistry) Register() error {
	that.mu.RLock()
	config := that.registry.Options()
	rsvc := that.service
	if config.Context == nil {
		that.mu.RUnlock()
		return errors.New("Not found service name! ")
	}
	var ok bool
	if that.serviceName == "" {
		that.serviceName, ok = config.Context.Value("ServiceName").(string)
		if !ok {
			that.mu.RUnlock()
			return errors.New("Not found service name! ")
		}
	}
	if that.serviceVersion == "" {
		that.serviceVersion, ok = config.Context.Value("ServiceVersion").(string)
		if !ok {
			that.mu.RUnlock()
			return errors.New("Not found service version! ")
		}
	}
	that.mu.RUnlock()

	regFunc := func(service *Service) error {
		ttl := DefaultRegisterTTL
		if config.Context != nil {
			ttl1, ok := config.Context.Value("RegisterTTL").(time.Duration)
			if ok {
				ttl = ttl1
			}
		}

		rOpts := []RegisterOption{OptRegisterTTL(ttl)}
		var regErr error

		for i := 0; i < 3; i++ {
			// 循环注册三次，直到注册成功，如果三次都为成功，则报错
			if e := that.registry.Register(service, rOpts...); e != nil {
				regErr = e
				time.Sleep(backoff.DoMul(i + 1))
				continue
			}
			// 如果执行三次，最后成功了，则清除错误
			regErr = nil
			break
		}

		return regErr
	}
	if rsvc != nil {
		if err := regFunc(rsvc); err != nil {
			return err
		}
		return nil
	}
	node := &Node{
		Id:       that.serviceName + "-" + that.addr.String(),
		Address:  that.addr.String(),
		Metadata: make(map[string]string),
		Paths:    that.allApis.Slice(),
	}
	node.Metadata["registry"] = that.registry.String()
	node.Metadata["server"] = that.serviceName

	svr := &Service{
		Name:    that.serviceName,
		Version: that.serviceVersion,
		Nodes:   []*Node{node},
	}
	err := regFunc(svr)
	if err != nil {
		return err
	}
	that.mu.Lock()
	that.registered = true
	that.service = svr
	that.mu.Unlock()
	return nil
}

// BeforeCloseEndpoint 关闭服务之前先取消注册
func (that *PluginRegistry) BeforeCloseEndpoint(_ drpc.Endpoint) error {
	err := that.Deregister()
	that.exit <- true
	return err
}

// Deregister 取消注册
func (that *PluginRegistry) Deregister() error {
	node := &Node{
		Id:      that.serviceName + "-" + that.addr.String(),
		Address: that.addr.String(),
	}

	service := &Service{
		Name:    that.serviceName,
		Version: that.serviceVersion,
		Nodes:   []*Node{node},
	}
	return that.registry.Deregister(service)
}

// 获取
func (that *PluginRegistry) leasesInterval() time.Duration {
	leasesIVal, ok := that.registry.Options().Context.Value(leasesInterval{}).(time.Duration)
	if ok && leasesIVal > 0 {
		return leasesIVal
	}
	return DefaultRegisterInterval
}
