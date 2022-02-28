package selector

import (
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/registry/cache"
	"github.com/osgochina/dmicro/registry/etcd"
	"time"
)

// 注册表选择器
type registrySelector struct {
	so Options
	rc cache.Cache
}

// Init 初始化
func (that *registrySelector) Init(opts ...Option) error {
	for _, o := range opts {
		o(&that.so)
	}
	that.rc.Stop()
	that.rc = that.newCache()

	return nil
}

// Options 获取配置
func (that *registrySelector) Options() Options {
	return that.so
}

// Select 选择服务
func (that *registrySelector) Select(service string, opts ...SelectOption) (Next, error) {
	sOpts := SelectOptions{
		Strategy: that.so.Strategy,
	}
	for _, opt := range opts {
		opt(&sOpts)
	}
	// 通过缓存获取到服务列表
	services, err := that.rc.GetService(service)
	if err != nil {
		if err == registry.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	// 过滤服务
	for _, filter := range sOpts.Filters {
		services = filter(services)
	}
	// 没有可用的服务
	if len(services) == 0 {
		return nil, ErrNoneAvailable
	}
	// 选出可用的服务
	return sOpts.Strategy(services), nil
}

// Mark 设置针对节点的成功或错误
func (that *registrySelector) Mark(service string, node *registry.Node, err error) {
}

// Reset 重置服务的状态
func (that *registrySelector) Reset(service string) {
}

// Close 关闭选择器
func (that *registrySelector) Close() error {
	that.rc.Stop()

	return nil
}

// 选择器名称
func (that *registrySelector) String() string {
	return "registry"
}

// 创建注册表信息缓存
func (that *registrySelector) newCache() cache.Cache {
	opts := make([]cache.Option, 0, 1)
	if that.so.Context != nil {
		if t, ok := that.so.Context.Value("selector_ttl").(time.Duration); ok {
			opts = append(opts, cache.WithTTL(t))
		}
	}
	return cache.New(that.so.Registry, opts...)
}

// NewSelector 创建选择器
func NewSelector(opts ...Option) Selector {
	sOpt := Options{
		Strategy: Random,
	}

	for _, opt := range opts {
		opt(&sOpt)
	}

	if sOpt.Registry == nil {
		sOpt.Registry = etcd.NewRegistry()
	}
	s := &registrySelector{
		so: sOpt,
	}
	s.rc = s.newCache()
	return s
}
