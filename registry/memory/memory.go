package memory

import (
	"context"
	"github.com/osgochina/dmicro/registry"
	"sync"
)

type memRegistry struct {
	options registry.Options
	sync.RWMutex
}

// NewRegistry 创建内存注册中心
func NewRegistry(opts ...registry.Option) registry.Registry {
	options := registry.Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}

	reg := &memRegistry{
		options: options,
	}

	return reg
}

func (that *memRegistry) Init(opts ...registry.Option) error {
	for _, o := range opts {
		o(&that.options)
	}

	return nil
}

func (that *memRegistry) Options() registry.Options {
	return that.options
}

func (that *memRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	return nil
}

func (that *memRegistry) Deregister(s *registry.Service, opts ...registry.DeregisterOption) error {

	return nil
}

func (that *memRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	return nil, nil
}

func (that *memRegistry) ListServices(opts ...registry.ListOption) ([]*registry.Service, error) {

	return nil, nil
}

func (that *memRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return nil, nil
}

func (that *memRegistry) String() string {
	return "memory"
}
