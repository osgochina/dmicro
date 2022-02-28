package cache

import (
	"github.com/osgochina/dmicro/registry"
	"math/rand"
	"sync"
	"time"
)

var (
	DefaultTTL = time.Minute // 默认生存时间
)

type Cache interface {
	registry.Registry
}

type cache struct {
	// 注册表对象
	registry.Registry
	// 配置参数
	opts Options
	// 锁
	sync.RWMutex
	// 是否退出缓存
	exit chan bool
}

func (c *cache) get(service string) ([]*registry.Service, error) {

	return nil, nil
}

// 判断服务是否可用
func (c *cache) isValid(services []*registry.Service, ttl time.Time) bool {

	// 不存在services
	if len(services) == 0 {
		return false
	}
	// 声明周期为0,表示每次都应该获取新的，不缓存
	if ttl.IsZero() {
		return false
	}
	// 如果已经结束了声明周期
	if time.Since(ttl) > 0 {
		return false
	}

	return true
}

func (c *cache) GetService(service string, opts ...registry.GetOption) ([]*registry.Service, error) {
	// 获取服务列表
	services, err := c.get(service)
	if err != nil {
		return nil, err
	}
	// 未找到服务
	if len(services) == 0 {
		return nil, registry.ErrNotFound
	}
	return services, nil
}

// Stop 停止缓存
func (c *cache) Stop() {
	c.Lock()
	defer c.Unlock()

	select {
	case <-c.exit:
		return
	default:
		close(c.exit)
	}
}

// 服务名
func (c *cache) String() string {
	return "cache"
}

// New 创建一个新的cache
func New(r registry.Registry, opts ...Option) Cache {
	rand.Seed(time.Now().UnixNano())
	options := Options{
		TTL: DefaultTTL,
	}
	for _, o := range opts {
		o(&options)
	}
	return &cache{
		Registry: r,
		opts:     options,
		exit:     make(chan bool),
	}
}
