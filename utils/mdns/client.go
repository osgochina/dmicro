package mdns

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/osgochina/dmicro/logger"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"strings"
	"sync"
	"time"
)

// ServiceEntry 查询后的条目
type ServiceEntry struct {
	Name       string // 服务名
	Host       string //
	AddrV4     net.IP
	AddrV6     net.IP
	Port       int
	Info       string
	InfoFields []string
	TTL        int
	Type       uint16
	Addr       net.IP
	hasTXT     bool
	sent       bool
}

// 检查必填信息是否存在
func (that *ServiceEntry) complete() bool {

	return (len(that.AddrV4) > 0 || len(that.AddrV6) > 0 || len(that.Addr) > 0) && that.Port != 0 && that.hasTXT
}

// QueryParam 自定义查询参数
type QueryParam struct {
	Service             string               // 要查找的服务
	Domain              string               // 要查找的域。默认是 local
	Type                uint16               // dns查询类型，默认是 dns.TypePTR
	Context             context.Context      // 上下文
	Timeout             time.Duration        // 查询超时时间，默认1s，如果通过context提供了超时，则忽略该参数
	Interface           *net.Interface       // 使用组播的网卡
	Entries             chan<- *ServiceEntry // 收到的响应包，以channel形式提供
	WantUniCastResponse bool                 // 是否需要单播响应，参考 per 5.4 in RFC
}

// DefaultParams 获取查询默认参数
func DefaultParams(service string) *QueryParam {
	return &QueryParam{
		Service:             service,
		Domain:              "local",
		Timeout:             time.Second,
		Entries:             make(chan *ServiceEntry),
		WantUniCastResponse: false, // TODO(reddaly): Change this default.
	}
}

// Query 在一个域中查找一个给定的服务，最多等待指定的秒数
// 结果是流式的 一个通道。发送不会阻塞，所以客户应该确保 读取或缓存。
func Query(params *QueryParam) error {
	cli, err := newClient()
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Close()
	}()
	if params.Interface != nil {
		if err = cli.setInterface(params.Interface, false); err != nil {
			return err
		}
	}
	if params.Domain == "" {
		params.Domain = "local"
	}
	if params.Context == nil {
		if params.Timeout == 0 {
			params.Timeout = time.Second
		}
		params.Context, _ = context.WithTimeout(context.Background(), params.Timeout)
		if err != nil {
			return err
		}
	}
	return cli.query(params)
}

// Listen 无限期的监听多播的更新
func Listen(entries chan<- *ServiceEntry, exit chan struct{}) error {
	cli, err := newClient()
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Close()
	}()

	_ = cli.setInterface(nil, true)

	msgCh := make(chan *dns.Msg, 32)

	go cli.recv(cli.ipv4UniCastConn, msgCh)
	go cli.recv(cli.ipv6UniCastConn, msgCh)
	go cli.recv(cli.ipv4MulticastConn, msgCh)
	go cli.recv(cli.ipv6MulticastConn, msgCh)

	ip := make(map[string]*ServiceEntry)
	for {
		select {
		case <-exit:
			return nil
		case <-cli.closedCh:
			return nil
		case m := <-msgCh:
			e := messageToEntry(m, ip)
			if e == nil {
				continue
			}

			// Check if this entry is complete
			if e.complete() {
				if e.sent {
					continue
				}
				e.sent = true
				entries <- e
				ip = make(map[string]*ServiceEntry)
			} else {
				// Fire off a node specific query
				msg := new(dns.Msg)
				msg.SetQuestion(e.Name, dns.TypePTR)
				msg.RecursionDesired = false
				if err = cli.sendQuery(msg); err != nil {
					logger.Printf("[ERR] mdns: Failed to query instance %s: %v", e.Name, err)
				}
			}
		}
	}
}

// Lookup 使用默认参数查询
func Lookup(service string, entries chan<- *ServiceEntry) error {
	params := DefaultParams(service)
	params.Entries = entries
	return Query(params)
}

// 提供了一个查询接口，可以用于 使用mDNS搜索服务提供商
type client struct {
	ipv4UniCastConn *net.UDPConn // 单播ipv4地址
	ipv6UniCastConn *net.UDPConn // 单播ipv6地址

	ipv4MulticastConn *net.UDPConn // 多播Ipv4地址
	ipv6MulticastConn *net.UDPConn // 多播ipv6地址

	closed    bool
	closedCh  chan struct{}
	closeLock sync.Mutex
}

// 建立客户端
func newClient() (*client, error) {

	uConn4, err4 := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	uConn6, err6 := net.ListenUDP("udp6", &net.UDPAddr{IP: net.IPv6zero, Port: 0})

	if err4 != nil && err6 != nil {
		logger.Printf("[ERR] mdns: Failed to bind to udp port: %v %v", err4, err6)
	}
	if uConn4 == nil && uConn6 == nil {
		return nil, fmt.Errorf("failed to bind to any unicast udp port")
	}

	if uConn4 == nil {
		uConn4 = &net.UDPConn{}
	}

	if uConn6 == nil {
		uConn6 = &net.UDPConn{}
	}

	mConn4, err4 := net.ListenUDP("udp4", mdnsWildcardAddrIPv4)
	mConn6, err6 := net.ListenUDP("udp6", mdnsWildcardAddrIPv6)
	if err4 != nil && err6 != nil {
		logger.Printf("[ERR] mdns: Failed to bind to udp port: %v %v", err4, err6)
	}

	if mConn4 == nil && mConn6 == nil {
		return nil, fmt.Errorf("failed to bind to any multicast udp port")
	}

	if mConn4 == nil {
		mConn4 = &net.UDPConn{}
	}

	if mConn6 == nil {
		mConn6 = &net.UDPConn{}
	}
	p1 := ipv4.NewPacketConn(mConn4)
	p2 := ipv6.NewPacketConn(mConn6)
	_ = p1.SetMulticastLoopback(true)
	_ = p2.SetMulticastLoopback(true)

	// 所有网卡的多播组都加入
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var errCount1, errCount2 int

	for _, iFace := range iFaces {
		if err = p1.JoinGroup(&iFace, &net.UDPAddr{IP: mdnsGroupIPv4}); err != nil {
			errCount1++
		}
		if err = p2.JoinGroup(&iFace, &net.UDPAddr{IP: mdnsGroupIPv6}); err != nil {
			errCount2++
		}
	}

	if len(iFaces) == errCount1 && len(iFaces) == errCount2 {
		return nil, fmt.Errorf("Failed to join multicast group on all interfaces! ")
	}
	c := &client{
		ipv4MulticastConn: mConn4,
		ipv6MulticastConn: mConn6,
		ipv4UniCastConn:   uConn4,
		ipv6UniCastConn:   uConn6,
		closedCh:          make(chan struct{}),
	}
	return c, nil
}

// Close 关闭客户端
func (that *client) Close() error {
	that.closeLock.Lock()
	defer that.closeLock.Unlock()

	if that.closed {
		return nil
	}
	that.closed = true

	close(that.closedCh)

	if that.ipv4UniCastConn != nil {
		_ = that.ipv4UniCastConn.Close()
	}
	if that.ipv6UniCastConn != nil {
		_ = that.ipv6UniCastConn.Close()
	}
	if that.ipv4MulticastConn != nil {
		_ = that.ipv4MulticastConn.Close()
	}
	if that.ipv6MulticastConn != nil {
		_ = that.ipv6MulticastConn.Close()
	}

	return nil
}

// 指定发送查询包的网卡
func (that *client) setInterface(iFace *net.Interface, loopBack bool) error {
	p := ipv4.NewPacketConn(that.ipv4UniCastConn)
	if err := p.JoinGroup(iFace, &net.UDPAddr{IP: mdnsGroupIPv4}); err != nil {
		return err
	}
	p2 := ipv6.NewPacketConn(that.ipv6UniCastConn)
	if err := p2.JoinGroup(iFace, &net.UDPAddr{IP: mdnsGroupIPv6}); err != nil {
		return err
	}
	p = ipv4.NewPacketConn(that.ipv4MulticastConn)
	if err := p.JoinGroup(iFace, &net.UDPAddr{IP: mdnsGroupIPv4}); err != nil {
		return err
	}
	p2 = ipv6.NewPacketConn(that.ipv6MulticastConn)
	if err := p2.JoinGroup(iFace, &net.UDPAddr{IP: mdnsGroupIPv6}); err != nil {
		return err
	}

	if loopBack {
		_ = p.SetMulticastLoopback(true)
		_ = p2.SetMulticastLoopback(true)
	}

	return nil
}

// Query用于执行查询和流结果
func (that *client) query(params *QueryParam) error {

	// 生成服务名
	serviceAddr := fmt.Sprintf("%s.%s.", trimDot(params.Service), trimDot(params.Domain))

	msgCh := make(chan *dns.Msg, 32)
	go that.recv(that.ipv4UniCastConn, msgCh)
	go that.recv(that.ipv6UniCastConn, msgCh)
	go that.recv(that.ipv4MulticastConn, msgCh)
	go that.recv(that.ipv6MulticastConn, msgCh)

	// 发送查询消息
	m := new(dns.Msg)
	if params.Type == dns.TypeNone {
		m.SetQuestion(serviceAddr, dns.TypePTR)
	} else {
		m.SetQuestion(serviceAddr, params.Type)
	}
	// RFC 6762, section 18.12.  Repurposing of Top Bit of qclass in Question
	// Section
	//
	// In the Question Section of a Multicast DNS query, the top bit of the qclass
	// field is used to indicate that unicast responses are preferred for this
	// particular question.  (See Section 5.4.)
	if params.WantUniCastResponse {
		m.Question[0].Qclass |= 1 << 15
	}
	m.RecursionDesired = false
	if err := that.sendQuery(m); err != nil {
		return err
	}

	inProgress := make(map[string]*ServiceEntry)
	for {
		select {
		case resp := <-msgCh:
			inp := messageToEntry(resp, inProgress)

			if inp == nil {
				continue
			}
			if len(resp.Question) == 0 || resp.Question[0].Name != m.Question[0].Name {
				// discard anything which we've not asked for
				continue
			}
			// Check if this entry is complete
			if inp.complete() {
				if inp.sent {
					continue
				}

				inp.sent = true
				select {
				case params.Entries <- inp:
				case <-params.Context.Done():
					return nil
				}
			} else {
				// Fire off a node specific query
				msg := new(dns.Msg)
				msg.SetQuestion(inp.Name, inp.Type)
				msg.RecursionDesired = false
				if err := that.sendQuery(msg); err != nil {
					logger.Printf("[ERR] mdns: Failed to query instance %s: %v", inp.Name, err)
				}
			}
		case <-params.Context.Done():
			return nil
		}
	}
}

// sendQuery 使用单播的方式发送dns查询包
func (that *client) sendQuery(q *dns.Msg) error {
	buf, err := q.Pack()
	if err != nil {
		return err
	}
	if that.ipv4UniCastConn != nil {
		_, _ = that.ipv4UniCastConn.WriteToUDP(buf, ipv4Addr)
	}
	if that.ipv6UniCastConn != nil {
		_, _ = that.ipv6UniCastConn.WriteToUDP(buf, ipv6Addr)
	}
	return nil
}

// 接收组播包
func (that *client) recv(l *net.UDPConn, msgCh chan *dns.Msg) {
	if l == nil {
		return
	}
	buf := make([]byte, 65536)
	for {
		that.closeLock.Lock()
		if that.closed {
			that.closeLock.Unlock()
			return
		}
		that.closeLock.Unlock()
		n, err := l.Read(buf)
		if err != nil {
			continue
		}
		msg := new(dns.Msg)
		if err = msg.Unpack(buf[:n]); err != nil {
			continue
		}
		select {
		case msgCh <- msg:
		case <-that.closedCh:
			return
		}
	}
}

// ensureName 转换
func ensureName(inProgress map[string]*ServiceEntry, name string, typ uint16) *ServiceEntry {
	if inp, ok := inProgress[name]; ok {
		return inp
	}
	inp := &ServiceEntry{
		Name: name,
		Type: typ,
	}
	inProgress[name] = inp
	return inp
}

// 别名
func alias(inProgress map[string]*ServiceEntry, src, dst string, typ uint16) {
	srcEntry := ensureName(inProgress, src, typ)
	inProgress[dst] = srcEntry
}

// dns消息转换为条目
func messageToEntry(m *dns.Msg, inProgress map[string]*ServiceEntry) *ServiceEntry {
	var inp *ServiceEntry

	for _, answer := range append(m.Answer, m.Extra...) {
		switch rr := answer.(type) {
		case *dns.PTR:
			// Create new entry for this
			inp = ensureName(inProgress, rr.Ptr, rr.Hdr.Rrtype)
			if inp.complete() {
				continue
			}
		case *dns.SRV:
			// Check for a target mismatch
			if rr.Target != rr.Hdr.Name {
				alias(inProgress, rr.Hdr.Name, rr.Target, rr.Hdr.Rrtype)
			}

			// Get the port
			inp = ensureName(inProgress, rr.Hdr.Name, rr.Hdr.Rrtype)
			if inp.complete() {
				continue
			}
			inp.Host = rr.Target
			inp.Port = int(rr.Port)
		case *dns.TXT:
			// Pull out the txt
			inp = ensureName(inProgress, rr.Hdr.Name, rr.Hdr.Rrtype)
			if inp.complete() {
				continue
			}
			inp.Info = strings.Join(rr.Txt, "|")
			inp.InfoFields = rr.Txt
			inp.hasTXT = true
		case *dns.A:
			// Pull out the IP
			inp = ensureName(inProgress, rr.Hdr.Name, rr.Hdr.Rrtype)
			if inp.complete() {
				continue
			}
			inp.Addr = rr.A // @Deprecated
			inp.AddrV4 = rr.A
		case *dns.AAAA:
			// Pull out the IP
			inp = ensureName(inProgress, rr.Hdr.Name, rr.Hdr.Rrtype)
			if inp.complete() {
				continue
			}
			inp.Addr = rr.AAAA // @Deprecated
			inp.AddrV6 = rr.AAAA
		}
		if inp != nil {
			inp.TTL = int(answer.Header().Ttl)
		}
	}

	return inp
}
