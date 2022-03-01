package cache

import (
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	registryUtil "github.com/osgochina/dmicro/registry/util"
	"github.com/osgochina/dmicro/utils/backoff"
	"golang.org/x/sync/singleflight"
	"math/rand"
	"sync"
	"time"
)

var (
	DefaultTTL = time.Minute // 默认生存时间
)

type Cache interface {
	registry.Registry
	Stop()
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

	// 服务缓存
	cache map[string][]*registry.Service
	// 服务缓存的生命周期
	ttls map[string]time.Time
	// 服务状态监听列表
	watched map[string]bool
	// 是否开启协程监听服务
	watchRunning bool
	// 防止缓存击穿
	sg singleflight.Group
	// 注册中心的故障状态
	status error
}

func (that *cache) get(service string) ([]*registry.Service, error) {
	// 增加读锁
	that.RLock()

	services := that.cache[service]
	ttl := that.ttls[service]

	cp := registryUtil.Copy(services)

	if that.isValid(cp, ttl) {
		that.RUnlock()
		return cp, nil
	}

	get := func(service string, cached []*registry.Service) ([]*registry.Service, error) {
		val, err, _ := that.sg.Do(service, func() (interface{}, error) {
			return that.Registry.GetService(service)
		})
		serviceList, _ := val.([]*registry.Service)
		if err != nil {
			if len(cached) > 0 {
				// 设置服务注册中心的故障状态
				that.setStatus(err)
				return cached, nil
			}
			return nil, err
		}
		// 如果上一次执行出现故障，因为此次成功，说明服务注册中心已经恢复，把错误设置为空
		if e := that.getStatus(); e != nil {
			that.setStatus(nil)
		}
		// 设置缓存
		that.Lock()
		that.set(service, registryUtil.Copy(serviceList))
		that.Unlock()
		return serviceList, nil
	}

	// 服务是否已经被监听？
	_, ok := that.watched[service]
	// 解开读锁
	that.RUnlock()

	// 如果该服务状态还未在监听中，则进入监听
	if !ok {
		that.Lock()
		that.watched[service] = true

		//判断是否开启了watch协程
		if !that.watchRunning {
			go that.run(service)
		}

		that.Unlock()
	}

	return get(service, services)
}

// 新开一个协程，监听服务
func (that *cache) run(service string) {
	that.Lock()
	that.watchRunning = true
	that.Unlock()

	// 该协程结束的时候，初始化监听列表，并且把运行状态设置为false
	defer func() {
		that.Lock()
		that.watched = make(map[string]bool)
		that.watchRunning = false
		that.Unlock()
	}()
	var a, b int
	for {
		//如果watch已经退出，则结束该协程
		if that.quit() {
			return
		}
		// 生成随机数，暂停随机毫秒数，防止并发冲突
		j := rand.Int63n(100)
		time.Sleep(time.Duration(j) * time.Millisecond)

		// 监听服务
		w, err := that.Registry.Watch(registry.WatchService(service))
		// 监听出错
		if err != nil {
			// 如果是监听服务结束了，则直接退出结束该协程
			if that.quit() {
				return
			}
			// 计算需要阻塞的时间
			d := backoff.Do(a)
			that.setStatus(err)
			// 如果大于三次，则重置
			if a > 3 {
				logger.Debugf("rcache: %v backing off %d", err, d)
				a = 0
			}
			time.Sleep(d)
			a++
			continue
		}
		// 重置重试次数
		a = 0
		err = that.watch(w)
		if err != nil {
			if that.quit() {
				return
			}
			//计算阻塞时间
			d := backoff.Do(b)
			that.setStatus(err)
			if b > 3 {
				logger.Debugf("rcache: %v backing off %d", err, d)
				b = 0
			}
			time.Sleep(d)
			b++
			continue
		}
		//重置重试次数
		b = 0
	}

}

// 判断服务是否可用
func (that *cache) isValid(services []*registry.Service, ttl time.Time) bool {

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

// 获取服务注册中心最后的故障状态
func (that *cache) getStatus() error {
	that.RLock()
	defer that.RUnlock()
	return that.status
}

// 设置服务注册中心的故障状态
func (that *cache) setStatus(err error) {
	that.Lock()
	that.status = err
	that.Unlock()
}

// 设置服务列表到缓存
func (that *cache) set(service string, services []*registry.Service) {
	that.cache[service] = services
	that.ttls[service] = time.Now().Add(that.opts.TTL)
}

// 判断当前watch监听是否已经退出
func (that *cache) quit() bool {
	select {
	case <-that.exit:
		return true
	default:
		return false
	}
}

// 监听服务配置
func (that *cache) watch(w registry.Watcher) error {
	stop := make(chan bool)

	go func() {
		defer w.Stop()
		select {
		case <-that.exit:
			return
		case <-stop:
			return
		}
	}()

	for {
		res, err := w.Next()
		if err != nil {
			close(stop)
			return err
		}
		// 重置服务注册中心的错误
		if e := that.getStatus(); e != nil {
			that.setStatus(nil)
		}
		that.update(res)
	}
}

// 更新配置信息到内存
func (that *cache) update(res *registry.Result) {
	if res == nil || res.Service == nil {
		return
	}
	that.Lock()
	defer that.Unlock()

	// 如果变更的服务不在监听列表，则忽略
	if _, ok := that.watched[res.Service.Name]; !ok {
		return
	}
	// 如果服务不在缓存列表，则说明缓存已经被删除了，不需要更新
	services, ok := that.cache[res.Service.Name]
	if !ok {
		return
	}
	// 如果服务的节点已经不存在了，且是删除配置的时间，则从缓存中删除该服务
	if len(res.Service.Nodes) == 0 {
		switch res.Action {
		case "delete":
			that.del(res.Service.Name)
		}
		return
	}
	// 查找缓存中存在的配置信息
	var service *registry.Service
	var index int
	for i, s := range services {
		if s.Version == res.Service.Version {
			service = s
			index = i
		}
	}
	switch res.Action {
	case "create", "update":
		// 如果service为nil，则此次触发的事件，更新的信息不存在缓存中，则表示需要把该信息写入缓存
		if service == nil {
			that.set(res.Service.Name, append(services, res.Service))
			return
		}
		for _, cur := range service.Nodes {
			var seen bool
			for _, node := range res.Service.Nodes {
				if cur.Id == node.Id {
					seen = true
					break
				}
			}
			// 如果node节点的配置信息不存在缓存中，则需要把缓存中的节点追加到配置中
			if !seen {
				res.Service.Nodes = append(res.Service.Nodes, cur)
			}
		}
		// 更新配置信息到缓存
		services[index] = res.Service
		that.set(res.Service.Name, services)
	case "delete":
		//如果已经删除，则不需要在执行了
		if service == nil {
			return
		}
		var nodes []*registry.Node
		for _, cur := range service.Nodes {
			var seen bool
			for _, del := range res.Service.Nodes {
				if del.Id == cur.Id {
					seen = true
					break
				}
			}
			// 判断缓存中的节点是否在要被删除的列表中，如果该节点要被删除，则跳过，如果该节点不在待删除的集合中，则把它重新赋值给nodes
			if !seen {
				nodes = append(nodes, cur)
			}
		}
		// 正常的节点，把它写入缓存
		if len(nodes) > 0 {
			service.Nodes = nodes
			services[index] = service
			that.set(service.Name, services)
			return
		}
		// TODO 暂时没明白这个逻辑的含义
		if len(services) == 1 {
			that.del(service.Name)
			return
		}
		//多版本服务
		var srvS []*registry.Service
		for _, s := range services {
			if s.Version != service.Version {
				srvS = append(srvS, s)
			}
		}
		// 保存服务的多版本配置信息
		that.set(service.Name, srvS)
	case "override":
		// 如果是覆盖事件，则删除该缓存
		if service == nil {
			return
		}
		that.del(service.Name)
	}
}

// 从缓存中删除配置信息
func (that *cache) del(service string) {
	// 如果服务注册中心出现了故障，则不能主动删除内存中的配置信息，这样能防止意外情况发生，兜底服务的可用性
	if err := that.status; err != nil {
		return
	}
	// 删除配置
	delete(that.cache, service)
	delete(that.ttls, service)
}

func (that *cache) GetService(service string, opts ...registry.GetOption) ([]*registry.Service, error) {
	// 获取服务列表
	services, err := that.get(service)
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
func (that *cache) Stop() {
	that.Lock()
	defer that.Unlock()

	select {
	case <-that.exit:
		return
	default:
		close(that.exit)
	}
}

// 服务名
func (that *cache) String() string {
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
		cache:    make(map[string][]*registry.Service),
		ttls:     make(map[string]time.Time),
		watched:  make(map[string]bool),
	}
}
