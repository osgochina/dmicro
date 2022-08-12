package memory

import (
	"context"
	"github.com/gogf/gf/util/guid"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/registry"
	"sync"
	"time"
)

var (
	// 节点存活时间检查
	ttlPruneTime = time.Second
	// 发送event超时时间
	sendEventTime = 10 * time.Millisecond
)

type memRegistry struct {
	sync.RWMutex
	// 配置信息
	options registry.Options
	// 记录 map[serverName]map[serverVersion]record
	records map[string]map[string]*record
	// 监听器
	watchers map[string]*memWatcher
}

// NewRegistry 创建内存注册中心
func NewRegistry(opts ...registry.Option) registry.Registry {
	options := registry.Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	records := getServiceRecords(options.Context)
	if records == nil {
		records = make(map[string]map[string]*record)
	}

	reg := &memRegistry{
		options: options,
		records: records,
	}

	go reg.ttlPrune()

	return reg
}

// Init 初始化配置
func (that *memRegistry) Init(opts ...registry.Option) error {
	for _, o := range opts {
		o(&that.options)
	}
	that.Lock()
	defer that.Unlock()

	records := getServiceRecords(that.options.Context)

	for name, rs := range records {
		// 如果当前记录集中不存在该服务名的记录,则多版本的全部赋值
		if _, ok := that.records[name]; !ok {
			that.records[name] = rs
			continue
		}
		//根据版本赋值
		for version, r := range rs {
			if _, ok := that.records[name][version]; !ok {
				that.records[name][version] = r
				continue
			}
		}
	}

	return nil
}

// Options 获取配置
func (that *memRegistry) Options() registry.Options {
	return that.options
}

// Register 注册服务
func (that *memRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	that.Lock()
	defer that.Unlock()

	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	// 创建service记录器
	if _, ok := that.records[s.Name]; !ok {
		that.records[s.Name] = make(map[string]*record)
	}

	// 通过service对象,生成record对象,并且区分版本记录
	r := serviceToRecord(s, options.TTL)
	if _, ok := that.records[s.Name][s.Version]; !ok {
		that.records[s.Name][s.Version] = r
		logger.Debugf("Registry added new service: %s, version: %s", s.Name, s.Version)
		// 发送事件
		go that.event(&registry.Result{Action: "update", Service: s})
		return nil
	}
	addedNodes := false
	// 把服务节点写入记录器
	for _, n := range s.Nodes {
		//如果该节点未注册,则执行注册操作
		if _, ok := that.records[s.Name][s.Version].Nodes[n.Id]; !ok {
			addedNodes = true
			metadata := make(map[string]string)
			for k, v := range n.Metadata {
				metadata[k] = v
			}
			// 创建节点
			that.records[s.Name][s.Version].Nodes[n.Id] = &node{
				Node: &registry.Node{
					Id:       n.Id,
					Address:  n.Address,
					Metadata: metadata,
				},
				TTL:      options.TTL,
				LastSeen: time.Now(),
			}
		}
	}
	//如果有新的节点增加,需要发送事件
	if addedNodes {
		logger.Debugf("Registry added new node to service: %s, version: %s", s.Name, s.Version)
		go that.event(&registry.Result{Action: "update", Service: s})
		return nil
	}

	// 刷新已存在的节点生存时间
	for _, n := range s.Nodes {
		logger.Debugf("Updated registration for service: %s, version: %s", s.Name, s.Version)
		that.records[s.Name][s.Version].Nodes[n.Id].TTL = options.TTL
		that.records[s.Name][s.Version].Nodes[n.Id].LastSeen = time.Now()
	}
	return nil
}

// Deregister 注销服务
func (that *memRegistry) Deregister(s *registry.Service, opts ...registry.DeregisterOption) error {

	that.Lock()
	defer that.Unlock()

	// 不存在服务
	if _, ok := that.records[s.Name]; !ok {
		return nil
	}
	//不存在该服务的版本
	if _, ok := that.records[s.Name][s.Version]; !ok {
		return nil
	}
	// 找到要注销的节点,从records中删除
	for _, n := range s.Nodes {
		if _, ok := that.records[s.Name][s.Version].Nodes[n.Id]; ok {
			logger.Debugf("Registry removed node from service: %s, version: %s", s.Name, s.Version)
			delete(that.records[s.Name][s.Version].Nodes, n.Id)
		}
	}
	//如果服务该版本的节点不存在则把该版本删除
	if len(that.records[s.Name][s.Version].Nodes) == 0 {
		logger.Debugf("Registry removed service: %s, version: %s", s.Name, s.Version)
		delete(that.records[s.Name], s.Version)
	}
	//服务版本都不存在,则删除服务
	if len(that.records[s.Name]) == 0 {
		logger.Debugf("Registry removed service: %s", s.Name)
		delete(that.records, s.Name)
	}
	//触发事件
	go that.event(&registry.Result{Action: "delete", Service: s})

	return nil
}

// GetService 通过service name获取service对象
func (that *memRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	that.Lock()
	defer that.Unlock()
	records, ok := that.records[name]
	if !ok {
		return nil, registry.ErrNotFound
	}
	services := make([]*registry.Service, len(records))
	i := 0
	for _, r := range records {
		services[i] = recordToService(r)
		i++
	}
	return services, nil
}

// ListServices 获取所有service对象
func (that *memRegistry) ListServices(opts ...registry.ListOption) ([]*registry.Service, error) {

	that.RLock()
	defer that.RUnlock()

	var services []*registry.Service
	for _, rs := range that.records {
		for _, r := range rs {
			services = append(services, recordToService(r))
		}
	}
	return services, nil
}

// Watch 监视
func (that *memRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	var watchOpts registry.WatchOptions
	for _, o := range opts {
		o(&watchOpts)
	}

	w := &memWatcher{
		exit:      make(chan bool),
		result:    make(chan *registry.Result),
		id:        guid.S(),
		watchOpts: watchOpts,
	}
	that.Lock()
	that.watchers[w.id] = w
	that.Unlock()

	return w, nil
}

func (that *memRegistry) String() string {
	return "memory"
}

// 节点生命期修整器
// 定时判断节点是否已经超过生命周期
func (that *memRegistry) ttlPrune() {

	prune := time.NewTicker(ttlPruneTime)
	defer prune.Stop()

	for {
		select {
		case <-prune.C:
			that.Lock()
			for name, records := range that.records {
				for version, r := range records {
					for id, n := range r.Nodes {
						//节点的最后可用时间到现在已经太长了,超过了ttl
						if n.TTL != 0 && time.Since(n.LastSeen) > n.TTL {
							logger.Debugf("Registry TTL expired for node %s of service %s", n.Id, name)
							delete(that.records[name][version].Nodes, id)
						}
					}
				}
			}
			that.Unlock()
		}
	}
}

// 触发事件
func (that *memRegistry) event(r *registry.Result) {

	that.RLock()
	watchers := make([]*memWatcher, 0, len(that.watchers))
	for _, w := range that.watchers {
		watchers = append(watchers, w)
	}
	that.RUnlock()

	for _, w := range watchers {
		select {
		case <-w.exit:
			that.Lock()
			delete(that.watchers, w.id)
			that.Unlock()
		default:
			select {
			case w.result <- r:
			case <-time.After(sendEventTime):
			}
		}
	}

}
