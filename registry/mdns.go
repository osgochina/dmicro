package registry

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/mdns"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mdnsDomain = "dmicro"
)

type mdnsTxt struct {
	Service  string
	Version  string
	Metadata map[string]string
}

type mdnsEntry struct {
	id   string
	node *mdns.Server
}

type mdnsRegistry struct {
	sync.Mutex
	opts   Options
	domain string

	services map[string][]*mdnsEntry
	mtx      sync.RWMutex

	watchers map[string]*mdnsWatcher

	// listener
	listener chan *mdns.ServiceEntry
}

type mdnsWatcher struct {
	id       string
	wo       WatchOptions
	ch       chan *mdns.ServiceEntry
	exit     chan struct{}
	domain   string
	registry *mdnsRegistry
}

func newRegistry(opts ...Option) Registry {
	options := Options{
		Context: context.Background(),
		Timeout: time.Second,
	}
	for _, o := range opts {
		o(&options)
	}
	domain := mdnsDomain

	d, ok := options.Context.Value("mdns.domain").(string)
	if ok {
		domain = d
	}

	return &mdnsRegistry{
		opts:     options,
		domain:   domain,
		services: make(map[string][]*mdnsEntry),
		watchers: make(map[string]*mdnsWatcher),
	}
}

func (that *mdnsRegistry) Init(opts ...Option) error {
	for _, o := range opts {
		o(&that.opts)
	}
	return nil
}

func (that *mdnsRegistry) Options() Options {
	return that.opts
}

// Register 注册服务节点
func (that *mdnsRegistry) Register(service *Service, _ ...RegisterOption) error {
	that.Lock()
	defer that.Unlock()

	entries, ok := that.services[service.Name]

	if !ok {
		s, err := mdns.NewServiceMDNS(service.Name, "_services", that.domain+".", "", 9999, []net.IP{net.ParseIP("0.0.0.0")}, nil)
		if err != nil {
			return err
		}

		svr, err := mdns.NewServer(&mdns.Config{Zone: &mdns.DNSSDService{ServiceMDNS: s}})
		if err != nil {
			return err
		}
		entries = append(entries, &mdnsEntry{id: "*", node: svr})
	}

	var gerr error
	for _, node := range service.Nodes {
		var seen bool
		var e *mdnsEntry
		//判断要添加的节点是否已经存在
		for _, entry := range entries {
			if node.Id == entry.id {
				seen = true
				e = entry
				break
			}
		}
		if seen {
			continue
		} else {
			e = &mdnsEntry{}
		}
		// 把节点数据编码
		txt, err := encode(&mdnsTxt{
			Service:  service.Name,
			Version:  service.Version,
			Metadata: node.Metadata,
		})
		if err != nil {
			gerr = err
			continue
		}
		// 分割地址端口
		host, pt, err := net.SplitHostPort(node.Address)
		if err != nil {
			gerr = err
			continue
		}
		port, _ := strconv.Atoi(pt)
		logger.Debugf(context.TODO(), "[mdns] registry create new service with ip: %s:%d for: %s", net.ParseIP(host).String(), port, host)

		//
		s, err := mdns.NewServiceMDNS(node.Id, service.Name, that.domain+".", "", port, []net.IP{net.ParseIP(host)}, txt)
		if err != nil {
			gerr = err
			continue
		}
		srv, err := mdns.NewServer(&mdns.Config{Zone: s, LocalhostChecking: true})
		if err != nil {
			gerr = err
			continue
		}
		e.id = node.Id
		e.node = srv
		entries = append(entries, e)
	}

	that.services[service.Name] = entries

	return gerr
}

// Deregister 取消注册节点
func (that *mdnsRegistry) Deregister(service *Service, _ ...DeregisterOption) error {
	that.Lock()
	defer that.Unlock()

	var newEntries []*mdnsEntry

	for _, entry := range that.services[service.Name] {
		var remove bool
		for _, node := range service.Nodes {
			if node.Id == entry.id {
				_ = entry.node.Shutdown()
				remove = true
				continue
			}
		}
		// 如果不需要移除，则重新加入
		if !remove {
			newEntries = append(newEntries, entry)
		}
	}
	//如果最后节点只剩下一个，并且是通配符*，则删除它
	if len(newEntries) == 1 && newEntries[0].id == "*" {
		_ = newEntries[0].node.Shutdown()
		delete(that.services, service.Name)
	} else {
		that.services[service.Name] = newEntries
	}

	return nil
}

// GetService 获取指定的服务信息
func (that *mdnsRegistry) GetService(service string, _ ...GetOption) ([]*Service, error) {
	serviceMap := make(map[string]*Service)
	entries := make(chan *mdns.ServiceEntry, 10)
	done := make(chan bool)

	p := mdns.DefaultParams(service)
	// 设置查询时间，因为是查询是异步执行的，所以实际查询时间就是设置的超时时间
	var cancel context.CancelFunc
	p.Context, cancel = context.WithTimeout(context.Background(), that.opts.Timeout)
	defer cancel()

	// set entries channel
	p.Entries = entries
	// set the domain
	p.Domain = that.domain

	// 等待查询请求的响应
	go func() {
		for {
			select {
			case e := <-entries:
				// 跳过列表记录
				if p.Service == "_services" {
					continue
				}
				// 不是当前域
				if p.Domain != that.domain {
					continue
				}
				// 已过期
				if e.TTL == 0 {
					continue
				}
				txt, err := decode(e.InfoFields)
				if err != nil {
					continue
				}
				// 收到信息不是当前要获取的service
				if txt.Service != service {
					continue
				}
				// 版本不存在则新建
				s, ok := serviceMap[txt.Version]
				if !ok {
					s = &Service{
						Name:    txt.Service,
						Version: txt.Version,
					}
				}
				addr := ""
				if len(e.AddrV4) > 0 {
					addr = net.JoinHostPort(e.AddrV4.String(), fmt.Sprint(e.Port))
				} else if len(e.AddrV6) > 0 {
					addr = net.JoinHostPort(e.AddrV6.String(), fmt.Sprint(e.Port))
				} else {
					logger.Debugf(context.TODO(), "[mdns]: invalid endpoint received: %v", e)
					continue
				}
				s.Nodes = append(s.Nodes, &Node{
					Id:       strings.TrimSuffix(e.Name, "."+p.Service+"."+p.Domain+"."),
					Address:  addr,
					Metadata: txt.Metadata,
				})
				serviceMap[txt.Version] = s
			case <-p.Context.Done():
				close(done)
				return
			}
		}
	}()

	// 发送查询请求
	if err := mdns.Query(p); err != nil {
		return nil, err
	}
	// 等待响应
	<-done
	services := make([]*Service, 0, len(serviceMap))
	for _, s := range serviceMap {
		services = append(services, s)
	}
	return services, nil
}

// ListServices 获取服务列表
func (that *mdnsRegistry) ListServices(_ ...ListOption) ([]*Service, error) {
	serviceMap := make(map[string]bool)
	entries := make(chan *mdns.ServiceEntry, 10)
	done := make(chan bool)
	p := mdns.DefaultParams("_services")
	// set context with timeout
	var cancel context.CancelFunc
	p.Context, cancel = context.WithTimeout(context.Background(), that.opts.Timeout)
	defer cancel()
	// set entries channel
	p.Entries = entries
	// set domain
	p.Domain = that.domain

	var services []*Service

	go func() {
		for {
			select {
			case e := <-entries:
				if e.TTL == 0 {
					continue
				}
				if !strings.HasSuffix(e.Name, p.Domain+".") {
					continue
				}
				name := strings.TrimSuffix(e.Name, "."+p.Service+"."+p.Domain+".")
				if !serviceMap[name] {
					serviceMap[name] = true
					services = append(services, &Service{Name: name})
				}
			case <-p.Context.Done():
				close(done)
				return
			}
		}
	}()
	// execute query
	if err := mdns.Query(p); err != nil {
		return nil, err
	}

	// wait till done
	<-done

	return services, nil
}

// Watch 监听
func (that *mdnsRegistry) Watch(opts ...WatchOption) (Watcher, error) {
	var wo WatchOptions
	for _, o := range opts {
		o(&wo)
	}
	md := &mdnsWatcher{
		id:       guid.S(),
		wo:       wo,
		ch:       make(chan *mdns.ServiceEntry, 32),
		exit:     make(chan struct{}),
		domain:   that.domain,
		registry: that,
	}
	that.mtx.Lock()
	defer that.mtx.Unlock()
	that.watchers[md.id] = md
	if that.listener != nil {
		return md, nil
	}

	go func() {
		for {
			that.mtx.Lock()
			// 如果监听列表为空，则不执行操作
			if len(that.watchers) == 0 {
				that.listener = nil
				that.mtx.Unlock()
				return
			}
			//如果已存在监听，则退出
			if that.listener != nil {
				that.mtx.Unlock()
				return
			}

			exit := make(chan struct{})
			ch := make(chan *mdns.ServiceEntry, 32)
			that.listener = ch

			that.mtx.Unlock()

			go func() {
				send := func(w *mdnsWatcher, e *mdns.ServiceEntry) {
					select {
					case w.ch <- e:
					default:
					}
				}
				for {
					select {
					case <-exit:
						return
					case e, ok := <-ch:
						if !ok {
							return
						}
						that.mtx.RLock()
						// 每次接收到了事件，都需要给所有的watcher发送一遍
						for _, w := range that.watchers {
							send(w, e)
						}
						that.mtx.RUnlock()
					}
				}
			}()
			// 阻塞启动监听
			_ = mdns.Listen(ch, exit)

			// 监听结束
			that.mtx.Lock()
			that.listener = nil
			close(ch)
			that.mtx.Unlock()
		}
	}()

	return md, nil
}

func (that *mdnsRegistry) String() string {
	return "mdns"
}

// Next 获取监听的节点变化
func (that *mdnsWatcher) Next() (*Result, error) {
	for {
		select {
		case e := <-that.ch:
			txt, err := decode(e.InfoFields)
			if err != nil {
				continue
			}
			if len(txt.Service) == 0 || len(txt.Version) == 0 {
				continue
			}
			// 如果服务不是要监听的服务，则跳过
			if len(that.wo.Service) > 0 && txt.Service != that.wo.Service {
				continue
			}
			var action EventType
			if e.TTL == 0 {
				action = Delete
			} else {
				action = Create
			}
			service := &Service{
				Name:    txt.Service,
				Version: txt.Version,
			}
			suffix := fmt.Sprintf(".%s.%s.", service.Name, that.domain)
			if !strings.HasSuffix(e.Name, suffix) {
				continue
			}
			var addr string
			if len(e.AddrV4) > 0 {
				addr = net.JoinHostPort(e.AddrV4.String(), fmt.Sprint(e.Port))
			} else if len(e.AddrV6) > 0 {
				addr = net.JoinHostPort(e.AddrV6.String(), fmt.Sprint(e.Port))
			} else {
				addr = e.Addr.String()
			}
			service.Nodes = append(service.Nodes, &Node{
				Id:       strings.TrimSuffix(e.Name, suffix),
				Address:  addr,
				Metadata: txt.Metadata,
			})
			return &Result{
				Action:  action,
				Service: service,
			}, nil
		case <-that.exit:
			return nil, ErrWatcherStopped
		}
	}
}
func (that *mdnsWatcher) Stop() {
	select {
	case <-that.exit:
		return
	default:
		close(that.exit)
		// remove self from the registry
		that.registry.mtx.Lock()
		delete(that.registry.watchers, that.id)
		that.registry.mtx.Unlock()
	}
}

func NewRegistry(opts ...Option) Registry {
	return newRegistry(opts...)
}

// 编码
func encode(txt *mdnsTxt) ([]string, error) {
	b, err := json.Marshal(txt)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	defer buf.Reset()

	w := zlib.NewWriter(&buf)
	if _, err := w.Write(b); err != nil {
		return nil, err
	}
	_ = w.Close()
	encoded := hex.EncodeToString(buf.Bytes())
	// individual txt limit
	if len(encoded) <= 255 {
		return []string{encoded}, nil
	}
	var record []string

	for len(encoded) > 255 {
		record = append(record, encoded[:255])
		encoded = encoded[255:]
	}

	record = append(record, encoded)

	return record, nil
}

func decode(record []string) (*mdnsTxt, error) {
	encoded := strings.Join(record, "")

	hr, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(hr)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return nil, err
	}

	rbuf, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}

	var txt *mdnsTxt

	if err := json.Unmarshal(rbuf, &txt); err != nil {
		return nil, err
	}

	return txt, nil
}
