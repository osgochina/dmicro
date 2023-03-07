package mdns

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/osgochina/dmicro/logger"
	"github.com/prometheus/common/log"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	mdnsGroupIPv4 = net.ParseIP("224.0.0.251")
	mdnsGroupIPv6 = net.ParseIP("ff02::fb")

	// mDNS wildcard addresses
	mdnsWildcardAddrIPv4 = &net.UDPAddr{
		IP:   net.ParseIP("224.0.0.0"),
		Port: 5353,
	}
	mdnsWildcardAddrIPv6 = &net.UDPAddr{
		IP:   net.ParseIP("ff02::"),
		Port: 5353,
	}
	// mDNS endpoint addresses
	ipv4Addr = &net.UDPAddr{
		IP:   mdnsGroupIPv4,
		Port: 5353,
	}
	ipv6Addr = &net.UDPAddr{
		IP:   mdnsGroupIPv6,
		Port: 5353,
	}
)

type Server struct {
	config       *Config
	ipv4List     *net.UDPConn
	ipv6List     *net.UDPConn
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
	wg           sync.WaitGroup
	outboundIP   net.IP
}

type GetMachineIP func() net.IP

type Config struct {
	Zone              Zone
	IFace             *net.Interface
	Port              int
	GetMachineIP      GetMachineIP
	LocalhostChecking bool
}

func NewServer(config *Config) (*Server, error) {
	setCustomPort(config.Port)

	// 监听指定的广播地址
	ipv4List, _ := net.ListenUDP("udp4", mdnsWildcardAddrIPv4)
	ipv6List, _ := net.ListenUDP("udp6", mdnsWildcardAddrIPv6)
	if ipv4List == nil && ipv6List == nil {
		return nil, fmt.Errorf("[ERR] mdns: Failed to bind to any udp port! ")
	}
	if ipv4List == nil {
		ipv4List = &net.UDPConn{}
	}
	if ipv6List == nil {
		ipv6List = &net.UDPConn{}
	}
	// 加入组播，接收通知
	p1 := ipv4.NewPacketConn(ipv4List)
	p2 := ipv6.NewPacketConn(ipv6List)
	_ = p1.SetMulticastLoopback(true)
	_ = p2.SetMulticastLoopback(true)
	// 如果指定的网卡不为空，则把该网卡的监听页加入
	if config.IFace != nil {
		if err := p1.JoinGroup(config.IFace, &net.UDPAddr{IP: mdnsGroupIPv4}); err != nil {
			return nil, err
		}
		if err := p2.JoinGroup(config.IFace, &net.UDPAddr{IP: mdnsGroupIPv6}); err != nil {
			return nil, err
		}
	} else {
		// 获取所有网卡列表
		ifaces, err := net.Interfaces()
		if err != nil {
			return nil, err
		}
		// 把每个网卡都加入
		errCount1, errCount2 := 0, 0
		for _, iface := range ifaces {
			if err := p1.JoinGroup(&iface, &net.UDPAddr{IP: mdnsGroupIPv4}); err != nil {
				errCount1++
			}
			if err := p2.JoinGroup(&iface, &net.UDPAddr{IP: mdnsGroupIPv6}); err != nil {
				errCount2++
			}
		}
		if len(ifaces) == errCount1 && len(ifaces) == errCount2 {
			return nil, fmt.Errorf("Failed to join multicast group on all interfaces! ")
		}
	}
	ipFunc := getOutboundIP
	if config.GetMachineIP != nil {
		ipFunc = config.GetMachineIP
	}

	s := &Server{
		config:     config,
		ipv4List:   ipv4List,
		ipv6List:   ipv6List,
		shutdownCh: make(chan struct{}),
		outboundIP: ipFunc(),
	}
	// 监听这些地址的收包
	go s.recv(s.ipv4List)
	go s.recv(s.ipv6List)
	s.wg.Add(1)
	go s.probe()
	return s, nil
}

func (that *Server) Shutdown() error {
	that.shutdownLock.Lock()
	defer that.shutdownLock.Unlock()
	if that.shutdown {
		return nil
	}
	that.shutdown = true
	close(that.shutdownCh)
	that.unregister()

	if that.ipv4List != nil {
		_ = that.ipv4List.Close()
	}
	if that.ipv6List != nil {
		_ = that.ipv6List.Close()
	}
	that.wg.Wait()
	return nil
}

// 设置自定义port
func setCustomPort(port int) {
	if port != 0 {
		if mdnsWildcardAddrIPv4.Port != port {
			mdnsWildcardAddrIPv4.Port = port
		}
		if mdnsWildcardAddrIPv6.Port != port {
			mdnsWildcardAddrIPv6.Port = port
		}
		if ipv4Addr.Port != port {
			ipv4Addr.Port = port
		}
		if ipv6Addr.Port != port {
			ipv6Addr.Port = port
		}
	}
}

// 获取这台机器的出口ip
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		// no net connectivity maybe so fallback
		return nil
	}
	defer func() {
		_ = conn.Close()
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// 接收广播包
func (that *Server) recv(c *net.UDPConn) {
	if c == nil {
		return
	}
	buf := make([]byte, 65536)
	for {
		// 如果服务已停止，则跳出
		that.shutdownLock.Lock()
		if that.shutdown {
			that.shutdownLock.Unlock()
			return
		}
		that.shutdownLock.Unlock()
		n, from, err := c.ReadFrom(buf)
		if err != nil {
			continue
		}
		// 解析服务收到的包
		if err := that.parsePacket(buf[:n], from); err != nil {
			log.Errorf("[ERR] mdns: Failed to handle query: %v", err)
		}
	}
}

// 解包
func (that *Server) parsePacket(packet []byte, from net.Addr) error {
	// 如果收到的包不是dns包，则报错
	var msg dns.Msg
	if err := msg.Unpack(packet); err != nil {
		log.Errorf("[ERR] mdns: Failed to unpack packet: %v", err)
		return err
	}
	// TODO: This is a bit of a hack
	// We decided to ignore some mDNS answers for the time being
	// See: https://tools.ietf.org/html/rfc6762#section-7.2
	msg.Truncated = false
	return that.handleQuery(&msg, from)
}

func (that *Server) handleQuery(query *dns.Msg, from net.Addr) error {
	//在组播查询和组播响应消息中，OPCODE必须在传输时为零(目前在多播中只支持标准查询)
	// 多播DNS消息接收与OPCODE其他 而零则必须被默默忽视。
	// 备注:OpcodeQuery == 0
	if query.Opcode != dns.OpcodeQuery {
		return fmt.Errorf("mdns: received query with non-zero Opcode %v: %v", query.Opcode, *query)
	}
	//在多播查询和多播响应消息中，响应传输时代码必须为零。
	//组播DNS报文接收非零响应码必须被安静地忽略。
	if query.Rcode != 0 {
		return fmt.Errorf("mdns: received query with non-zero Rcode %v: %v", query.Rcode, *query)
	}
	// TODO(reddaly): Handle "TC (Truncated) Bit":
	//    In query messages, if the TC bit is set, it means that additional
	//    Known-Answer records may be following shortly.  A responder SHOULD
	//    record this fact, and wait for those additional Known-Answer records,
	//    before deciding whether to respond.  If the TC bit is clear, it means
	//    that the querying host has no additional Known Answers.
	if query.Truncated {
		return fmt.Errorf("[ERR] mdns: support for DNS requests with high truncated bit not implemented: %v", *query)
	}

	var unicastAnswer, multicastAnswer []dns.RR
	// Handle each question
	for _, q := range query.Question {
		mrecs, urecs := that.handleQuestion(q)
		multicastAnswer = append(multicastAnswer, mrecs...)
		unicastAnswer = append(unicastAnswer, urecs...)
	}

	// See section 18 of RFC 6762 for rules about DNS headers.
	resp := func(unicast bool) *dns.Msg {
		// 18.1: ID (Query Identifier)
		// 0 for multicast response, query.Id for unicast response
		id := uint16(0)
		if unicast {
			id = query.Id
		}

		var answer []dns.RR
		if unicast {
			answer = unicastAnswer
		} else {
			answer = multicastAnswer
		}
		if len(answer) == 0 {
			return nil
		}

		return &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Id: id,

				// 18.2: QR (Query/Response) Bit - must be set to 1 in response.
				Response: true,

				// 18.3: OPCODE - must be zero in response (OpcodeQuery == 0)
				Opcode: dns.OpcodeQuery,

				// 18.4: AA (Authoritative Answer) Bit - must be set to 1
				Authoritative: true,

				// The following fields must all be set to 0:
				// 18.5: TC (TRUNCATED) Bit
				// 18.6: RD (Recursion Desired) Bit
				// 18.7: RA (Recursion Available) Bit
				// 18.8: Z (Zero) Bit
				// 18.9: AD (Authentic Data) Bit
				// 18.10: CD (Checking Disabled) Bit
				// 18.11: RCODE (Response Code)
			},
			// 18.12 pertains to questions (handled by handleQuestion)
			// 18.13 pertains to resource records (handled by handleQuestion)

			// 18.14: Name Compression - responses should be compressed (though see
			// caveats in the RFC), so set the Compress bit (part of the dns library
			// API, not part of the DNS packet) to true.
			Compress: true,
			Question: query.Question,
			Answer:   answer,
		}
	}
	if mresp := resp(false); mresp != nil {
		if err := that.sendResponse(mresp, from); err != nil {
			return fmt.Errorf("mdns: error sending multicast response: %v", err)
		}
	}
	if uresp := resp(true); uresp != nil {
		if err := that.sendResponse(uresp, from); err != nil {
			return fmt.Errorf("mdns: error sending unicast response: %v", err)
		}
	}

	return nil
}

func (that *Server) handleQuestion(q dns.Question) (multicastRecs, unicastRecs []dns.RR) {
	records := that.config.Zone.Records(q)
	if len(records) == 0 {
		return nil, nil
	}

	// Handle unicast and multicast responses.
	// TODO(reddaly): The decision about sending over unicast vs. multicast is not
	// yet fully compliant with RFC 6762.  For example, the unicast bit should be
	// ignored if the records in question are close to TTL expiration.  For now,
	// we just use the unicast bit to make the decision, as per the spec:
	//     RFC 6762, section 18.12.  Repurposing of Top Bit of qclass in Question
	//     Section
	//
	//     In the Question Section of a Multicast DNS query, the top bit of the
	//     qclass field is used to indicate that unicast responses are preferred
	//     for this particular question.  (See Section 5.4.)
	if q.Qclass&(1<<15) != 0 {
		return nil, records
	}
	return records, nil
}

// 发送响应消息
func (that *Server) sendResponse(resp *dns.Msg, from net.Addr) error {
	// TODO(reddaly): Respect the unicast argument, and allow sending responses
	// over multicast.
	buf, err := resp.Pack()
	if err != nil {
		return err
	}
	// Determine the socket to send from
	addr := from.(*net.UDPAddr)
	conn := that.ipv4List
	backupTarget := net.IPv4zero

	if addr.IP.To4() == nil {
		conn = that.ipv6List
		backupTarget = net.IPv6zero
	}
	_, err = conn.WriteToUDP(buf, addr)
	if that.config.LocalhostChecking && addr.IP.Equal(that.outboundIP) {
		// ignore any errors, this is best efforts
		_, _ = conn.WriteToUDP(buf, &net.UDPAddr{IP: backupTarget, Port: addr.Port})
	}
	return err
}

// 探针
func (that *Server) probe() {
	defer that.wg.Done()
	sd, ok := that.config.Zone.(*ServiceMDNS)
	if !ok {
		return
	}
	name := fmt.Sprintf("%s.%s.%s.", sd.Instance, trimDot(sd.Service), trimDot(sd.Domain))

	q := new(dns.Msg)
	q.SetQuestion(name, dns.TypePTR)
	q.RecursionDesired = false

	srv := &dns.SRV{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeSRV,
			Class:  dns.ClassINET,
			Ttl:    defaultTTL,
		},
		Priority: 0,
		Weight:   0,
		Port:     uint16(sd.Port),
		Target:   sd.HostName,
	}
	txt := &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    defaultTTL,
		},
		Txt: sd.TXT,
	}
	q.Ns = []dns.RR{srv, txt}
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 3; i++ {
		if err := that.SendMulticast(q); err != nil {
			logger.Errorf(context.TODO(), "[ERR] mdns: failed to send probe:%v", err.Error())
		}
		time.Sleep(time.Duration(randomizer.Intn(250)) * time.Millisecond)
	}
	resp := new(dns.Msg)
	resp.MsgHdr.Response = true
	// set for query
	q.SetQuestion(name, dns.TypeANY)
	resp.Answer = append(resp.Answer, that.config.Zone.Records(q.Question[0])...)
	// reset
	q.SetQuestion(name, dns.TypePTR)
	timeout := 1 * time.Second
	timer := time.NewTimer(timeout)
	for i := 0; i < 3; i++ {
		if err := that.SendMulticast(resp); err != nil {
			logger.Errorf(context.TODO(), "[ERR] mdns: failed to send announcement:%v", err.Error())
		}
		select {
		case <-timer.C:
			timeout *= 2
			timer.Reset(timeout)
		case <-that.shutdownCh:
			timer.Stop()
			return
		}
	}
}

// SendMulticast 发送组播消息
func (that *Server) SendMulticast(msg *dns.Msg) error {
	buf, err := msg.Pack()
	if err != nil {
		return err
	}
	if that.ipv4List != nil {
		_, _ = that.ipv4List.WriteToUDP(buf, ipv4Addr)
	}
	if that.ipv6List != nil {
		_, _ = that.ipv6List.WriteToUDP(buf, ipv6Addr)
	}
	return nil
}

// 发送广播消息，服务将要退出
func (that *Server) unregister() error {
	sd, ok := that.config.Zone.(*ServiceMDNS)
	if !ok {
		return nil
	}
	atomic.StoreUint32(&sd.TTL, 0)
	name := fmt.Sprintf("%s.%s.%s.", sd.Instance, trimDot(sd.Service), trimDot(sd.Domain))
	q := new(dns.Msg)
	q.SetQuestion(name, dns.TypeANY)

	resp := new(dns.Msg)
	resp.MsgHdr.Response = true
	resp.Answer = append(resp.Answer, that.config.Zone.Records(q.Question[0])...)

	return that.SendMulticast(resp)
}
