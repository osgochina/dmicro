package etcd

import (
	"context"
	"crypto/tls"
	"github.com/gogf/gf/encoding/ghash"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/osgochina/dmicro/registry"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	prefix = "/dmicro/registry/"
)

// etcd注册表
type etcdRegistry struct {
	client         *clientv3.Client            // etcd客户端
	options        registry.Options            // 参数
	register       map[string]uint64           // 已注册list
	leases         map[string]clientv3.LeaseID // 租约
	leasesInterval time.Duration               //租约续期时间
	sync.RWMutex
}

// NewRegistry 创建etcd注册表
func NewRegistry(opts ...registry.Option) registry.Registry {
	e := &etcdRegistry{
		options:  registry.Options{},
		register: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
	}
	// 如果有etcd的用户名密码环境变量
	username, password := os.Getenv("ETCD_USERNAME"), os.Getenv("ETCD_PASSWORD")
	if len(username) > 0 && len(password) > 0 {
		opts = append(opts, Auth(username, password))
	}
	// 如果有注册表地址环境变量
	address := os.Getenv("DMICRO_REGISTRY_ADDRESS")
	if len(address) > 0 {
		opts = append(opts, registry.AddrList(address))
	}
	// 执行配置解析
	_ = e.configure(opts...)
	return e
}

// 执行配置解析
func (that *etcdRegistry) configure(opts ...registry.Option) error {
	config := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}
	for _, o := range opts {
		o(&that.options)
	}
	if that.options.Timeout == 0 {
		that.options.Timeout = 5 * time.Second
	}
	config.DialTimeout = that.options.Timeout
	if that.options.Secure || that.options.TLSConfig != nil {
		tlsConfig := that.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		config.TLS = tlsConfig
	}
	if that.options.Context != nil {
		u, ok := that.options.Context.Value(authKey{}).(*authCreds)
		if ok {
			config.Username = u.Username
			config.Password = u.Password
		}
		cfg, ok := that.options.Context.Value(logConfigKey{}).(*zap.Config)
		if ok && cfg != nil {
			config.LogConfig = cfg
		}

	}
	var cAddrList []string

	for _, address := range that.options.AddrList {
		if len(address) == 0 {
			continue
		}
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "2379"
			addr = address
			cAddrList = append(cAddrList, net.JoinHostPort(addr, port))
		} else if err == nil {
			cAddrList = append(cAddrList, net.JoinHostPort(addr, port))
		}
	}
	if len(cAddrList) > 0 {
		config.Endpoints = cAddrList
	}
	cli, err := clientv3.New(config)
	if err != nil {
		return err
	}
	that.client = cli
	return nil
}

// Init 初始化参数
func (that *etcdRegistry) Init(opts ...registry.Option) error {
	return that.configure(opts...)
}

// Options 获取注册表的参数
func (that *etcdRegistry) Options() registry.Options {
	return that.options
}

// Register 注册服务
func (that *etcdRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if len(s.Nodes) == 0 {
		return gerror.New("至少要传入一个node节点")
	}
	var gerr error

	// 注册每一个节点
	for _, node := range s.Nodes {
		err := that.registerNode(s, node, opts...)
		if err != nil {
			gerr = err
		}
	}
	return gerr
}

// 注册节点
func (that *etcdRegistry) registerNode(s *registry.Service, node *registry.Node, opts ...registry.RegisterOption) (err error) {
	if len(s.Nodes) == 0 {
		return gerror.New("至少要传入一个node节点")
	}
	that.RLock()
	leaseID, ok := that.leases[s.Name+node.Id]
	that.RUnlock()

	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), that.options.Timeout)
		defer cancel()
		rsp, err := that.client.Get(ctx, nodePath(s.Name, node.Id), clientv3.WithSerializable())
		if err != nil {
			return err
		}
		for _, kv := range rsp.Kvs {
			if kv.Lease > 0 {
				leaseID = clientv3.LeaseID(kv.Lease)

				// 解码已存在的节点
				srv := decode(kv.Value)
				if srv == nil || len(srv.Nodes) == 0 {
					continue
				}

				// create hash of service; uint64
				h := ghash.BKDRHash64(gjson.New(srv.Nodes[0]).MustToJson())
				if err != nil {
					continue
				}

				// save the info
				that.Lock()
				that.leases[s.Name+node.Id] = leaseID
				that.register[s.Name+node.Id] = h
				that.Unlock()

				break
			}
		}
	}
	var leaseNotFound bool

	if leaseID > 0 {
		if _, err := that.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
			if err != rpctypes.ErrLeaseNotFound {
				return err
			}
			// 没有找到租约
			leaseNotFound = true
		}
	}
	h := ghash.BKDRHash64(gjson.New(node).MustToJson())
	that.Lock()
	v, ok := that.register[s.Name+node.Id]
	that.Unlock()
	// 如果该节点的租约已存在，则什么都不做
	if ok && v == h && !leaseNotFound {
		return nil
	}
	// 注册服务
	service := &registry.Service{
		Name:     s.Name,
		Version:  s.Version,
		Metadata: s.Metadata,
		Nodes:    []*registry.Node{node},
	}
	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithTimeout(context.Background(), that.options.Timeout)
	defer cancel()

	var lgr *clientv3.LeaseGrantResponse
	if options.TTL.Seconds() > 0 {
		// 如果设置了生存时间，则表示需要创建租约
		lgr, err = that.client.Grant(ctx, int64(options.TTL.Seconds()))
		if err != nil {
			return err
		}
	}
	// 创建了租约，则写入带租约的内容
	if lgr != nil {
		_, err = that.client.Put(ctx, nodePath(service.Name, node.Id), encode(service), clientv3.WithLease(lgr.ID))
	} else {
		_, err = that.client.Put(ctx, nodePath(service.Name, node.Id), encode(service))
	}
	if err != nil {
		return err
	}
	that.Lock()
	// 保存服务节点的hash值，用于下次比较
	that.register[s.Name+node.Id] = h
	// 保存服务节点的租约id
	if lgr != nil {
		that.leases[s.Name+node.Id] = lgr.ID
	}
	that.Unlock()

	return nil
}

// Deregister 取消注册
func (that *etcdRegistry) Deregister(s *registry.Service, _ ...registry.DeregisterOption) error {
	if len(s.Nodes) == 0 {
		return gerror.New("至少要传入一个node节点")
	}
	for _, node := range s.Nodes {
		that.Lock()
		// 从列表中删除服务节点的hash值
		delete(that.register, s.Name+node.Id)
		// 从租约列表中删除服务节点的租约id
		delete(that.leases, s.Name+node.Id)
		that.Unlock()

		// 从etcd的注册表中删除指定的节点信息
		ctx, cancel := context.WithTimeout(context.Background(), that.options.Timeout)
		defer cancel()

		_, err := that.client.Delete(ctx, nodePath(s.Name, node.Id))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetService 从etcd注册表中获取当前注册的节点列表，注意因为一个服务可以同时存在多个版本，所以需要按版本区分
func (that *etcdRegistry) GetService(name string, _ ...registry.GetOption) ([]*registry.Service, error) {

	// 获取当前注册的内容
	ctx, cancel := context.WithTimeout(context.Background(), that.options.Timeout)
	defer cancel()
	rsp, err := that.client.Get(ctx, servicePath(name)+"/", clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	// 如果该目录下未有值，则表示未注册
	if len(rsp.Kvs) == 0 {
		return nil, registry.ErrNotFound
	}
	serviceMap := map[string]*registry.Service{}
	// 按服务节点的版本区分服务
	for _, n := range rsp.Kvs {
		if sn := decode(n.Value); sn != nil {
			s, ok := serviceMap[sn.Version]
			if !ok {
				s = &registry.Service{
					Name:     sn.Name,
					Version:  sn.Version,
					Metadata: sn.Metadata,
				}
				serviceMap[s.Version] = s
			}
			s.Nodes = append(s.Nodes, sn.Nodes...)
		}
	}
	// 获取最后该名字的服务列表
	services := make([]*registry.Service, 0, len(serviceMap))
	for _, service := range serviceMap {
		services = append(services, service)
	}

	return services, nil
}

// ListServices 获取服务列表
func (that *etcdRegistry) ListServices(_ ...registry.ListOption) ([]*registry.Service, error) {
	// 根据版本区分服务列表
	versions := make(map[string]*registry.Service)

	ctx, cancel := context.WithTimeout(context.Background(), that.options.Timeout)
	defer cancel()

	// 获取该路径前缀下的所有值
	rsp, err := that.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	//如果不存在，则返回空
	if len(rsp.Kvs) == 0 {
		return []*registry.Service{}, nil
	}

	for _, n := range rsp.Kvs {
		sn := decode(n.Value)
		if sn == nil {
			continue
		}
		v, ok := versions[sn.Name+sn.Version]
		if !ok {
			versions[sn.Name+sn.Version] = sn
			continue
		}
		// 按服务名+版本区分节点
		v.Nodes = append(v.Nodes, sn.Nodes...)
	}
	// 把不同版本的服务组成数组
	services := make([]*registry.Service, 0, len(versions))
	for _, service := range versions {
		services = append(services, service)
	}

	// 按服务名排序服务列表
	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })

	return services, nil
}

// Watch 监听服务变更
func (that *etcdRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return newEtcdWatcher(that, that.options.Timeout, opts...)
}

func (that *etcdRegistry) String() string {
	return "etcd"
}
