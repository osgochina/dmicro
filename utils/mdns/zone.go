package mdns

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"os"
	"strings"
	"sync/atomic"
)

const (
	// defaultTTL 默认dns记录的存活时间
	defaultTTL = 120
)

type Zone interface {
	// Records 根据dns查询请求，返回一个dns记录。
	Records(q dns.Question) []dns.RR
}

// ServiceMDNS mdns服务
type ServiceMDNS struct {
	Instance     string   // 实例的名称 (e.g. "hostService name")
	Service      string   // 服务的名字 (e.g. "_http._tcp.")
	Domain       string   // 域 默认为“local”
	HostName     string   // 主机的dns名称 (e.g. "mymachine.net.")
	Port         int      // 端口
	IPs          []net.IP // 当前服务使用的ip地址
	TXT          []string
	TTL          uint32
	serviceAddr  string
	instanceAddr string
	enumAddr     string
}

// NewServiceMDNS 创建MDNS服务
// instance
func NewServiceMDNS(instance, service, domain, hostName string, port int, ips []net.IP, txt []string) (*ServiceMDNS, error) {
	if instance == "" {
		return nil, fmt.Errorf("missing service instance name")
	}
	if service == "" {
		return nil, fmt.Errorf("missing service name")
	}
	if port == 0 {
		return nil, fmt.Errorf("missing service port")
	}
	if domain == "" {
		domain = "local."
	}
	if err := validateFQDN(domain); err != nil {
		return nil, fmt.Errorf("domain %q is not a fully-qualified domain name: %v", domain, err)
	}
	if hostName == "" {
		var err error
		hostName, err = os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("could not determine host: %v", err)
		}
		hostName = fmt.Sprintf("%s.", hostName)
	}
	if err := validateFQDN(hostName); err != nil {
		return nil, fmt.Errorf("hostName %q is not a fully-qualified domain name: %v", hostName, err)
	}
	if len(ips) == 0 {
		var err error
		ips, err = net.LookupIP(trimDot(hostName))
		if err != nil {
			// Try appending the host domain suffix and lookup again
			// (required for Linux-based hosts)
			tmpHostName := fmt.Sprintf("%s%s", hostName, domain)

			ips, err = net.LookupIP(trimDot(tmpHostName))

			if err != nil {
				return nil, fmt.Errorf("could not determine host IP addresses for %s", hostName)
			}
		}
	}
	for _, ip := range ips {
		if ip.To4() == nil && ip.To16() == nil {
			return nil, fmt.Errorf("invalid IP address in IPs list: %v", ip)
		}
	}
	return &ServiceMDNS{
		Instance:     instance,
		Service:      service,
		Domain:       domain,
		HostName:     hostName,
		Port:         port,
		IPs:          ips,
		TXT:          txt,
		TTL:          defaultTTL,
		serviceAddr:  fmt.Sprintf("%s.%s.", trimDot(service), trimDot(domain)),
		instanceAddr: fmt.Sprintf("%s.%s.%s.", instance, trimDot(service), trimDot(domain)),
		enumAddr:     fmt.Sprintf("_services._dns-sd._udp.%s.", trimDot(domain)),
	}, nil
}

func validateFQDN(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("FQDN must not be blank")
	}
	if s[len(s)-1] != '.' {
		return fmt.Errorf("FQDN must end in period: %s", s)
	}
	// TODO(reddaly): Perform full validation.

	return nil
}

// 去除.号
func trimDot(s string) string {
	return strings.Trim(s, ".")
}

func (that *ServiceMDNS) Records(q dns.Question) []dns.RR {
	switch q.Name {
	case that.enumAddr:
		return that.serviceEnum(q)
	case that.serviceAddr:
		return that.serviceRecords(q)
	case that.instanceAddr:
		return that.instanceRecords(q)
	case that.HostName:
		if q.Qtype == dns.TypeA || q.Qtype == dns.TypeAAAA {
			return that.instanceRecords(q)
		}
		fallthrough
	default:
		return nil
	}
}

func (that *ServiceMDNS) serviceEnum(q dns.Question) []dns.RR {
	switch q.Qtype {
	case dns.TypeANY:
		fallthrough
	case dns.TypePTR:
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    atomic.LoadUint32(&that.TTL),
			},
			Ptr: that.serviceAddr,
		}
		return []dns.RR{rr}
	default:
		return nil
	}
}

func (that *ServiceMDNS) serviceRecords(q dns.Question) []dns.RR {
	switch q.Qtype {
	case dns.TypeANY:
		fallthrough
	case dns.TypePTR:
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    atomic.LoadUint32(&that.TTL),
			},
			Ptr: that.instanceAddr,
		}
		servRec := []dns.RR{rr}
		instRecs := that.instanceRecords(dns.Question{
			Name:  that.instanceAddr,
			Qtype: dns.TypeANY,
		})
		return append(servRec, instRecs...)
	default:
		return nil
	}
}

func (that *ServiceMDNS) instanceRecords(q dns.Question) []dns.RR {
	switch q.Qtype {
	case dns.TypeANY:
		// Get the SRV, which includes A and AAAA
		recs := that.instanceRecords(dns.Question{
			Name:  that.instanceAddr,
			Qtype: dns.TypeSRV,
		})

		// Add the TXT record
		recs = append(recs, that.instanceRecords(dns.Question{
			Name:  that.instanceAddr,
			Qtype: dns.TypeTXT,
		})...)
		return recs

	case dns.TypeA:
		var rr []dns.RR
		for _, ip := range that.IPs {
			if ip4 := ip.To4(); ip4 != nil {
				rr = append(rr, &dns.A{
					Hdr: dns.RR_Header{
						Name:   that.HostName,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    atomic.LoadUint32(&that.TTL),
					},
					A: ip4,
				})
			}
		}
		return rr

	case dns.TypeAAAA:
		var rr []dns.RR
		for _, ip := range that.IPs {
			if ip.To4() != nil {
				// TODO(reddaly): IPv4 addresses could be encoded in IPv6 format and
				// putinto AAAA records, but the current logic puts ipv4-encodable
				// addresses into the A records exclusively.  Perhaps this should be
				// configurable?
				continue
			}

			if ip16 := ip.To16(); ip16 != nil {
				rr = append(rr, &dns.AAAA{
					Hdr: dns.RR_Header{
						Name:   that.HostName,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    atomic.LoadUint32(&that.TTL),
					},
					AAAA: ip16,
				})
			}
		}
		return rr

	case dns.TypeSRV:
		// Create the SRV Record
		srv := &dns.SRV{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeSRV,
				Class:  dns.ClassINET,
				Ttl:    atomic.LoadUint32(&that.TTL),
			},
			Priority: 10,
			Weight:   1,
			Port:     uint16(that.Port),
			Target:   that.HostName,
		}
		recs := []dns.RR{srv}

		// Add the A record
		recs = append(recs, that.instanceRecords(dns.Question{
			Name:  that.instanceAddr,
			Qtype: dns.TypeA,
		})...)

		// Add the AAAA record
		recs = append(recs, that.instanceRecords(dns.Question{
			Name:  that.instanceAddr,
			Qtype: dns.TypeAAAA,
		})...)
		return recs

	case dns.TypeTXT:
		txt := &dns.TXT{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    atomic.LoadUint32(&that.TTL),
			},
			Txt: that.TXT,
		}
		return []dns.RR{txt}
	}
	return nil
}
