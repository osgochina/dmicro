package mdns

import "github.com/miekg/dns"

type DNSSDService struct {
	ServiceMDNS *ServiceMDNS
}

func (that *DNSSDService) Records(q dns.Question) []dns.RR {
	var recs []dns.RR
	if q.Name == "_services._dns-sd._udp."+that.ServiceMDNS.Domain+"." {
		recs = that.dnssdMetaQueryRecords(q)
	}
	return append(recs, that.ServiceMDNS.Records(q)...)
}

func (that *DNSSDService) dnssdMetaQueryRecords(q dns.Question) []dns.RR {
	return []dns.RR{
		&dns.PTR{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    defaultTTL,
			},
			Ptr: that.ServiceMDNS.serviceAddr,
		},
	}
}
