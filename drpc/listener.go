package drpc

import (
	"crypto/tls"
	"errors"
	"github.com/osgochina/dmicro/utils/inherit"
	"net"
)

// NewInheritedListener 创建一个支持优雅重启，支持继承监听的监听器
func NewInheritedListener(addr net.Addr, tlsConfig *tls.Config) (lis net.Listener, err error) {
	addrStr := addr.String()
	network := addr.Network()
	var host, port string
	switch addrF := addr.(type) {
	case *FakeAddr:
		host, port = addrF.Host(), addrF.Port()
	default:
		host, port, err = net.SplitHostPort(addrStr)
		if err != nil {
			return nil, err
		}
	}
	if port == "0" {
		addrStr = PopParentAddr(network, host, addrStr)
	}

	//if _network := asQUIC(network); _network != "" {
	//	if tlsConfig == nil {
	//		tlsConfig = testTLSConfig
	//	}
	//	lis, err = quic.InheritedListen(_network, laddr, tlsConfig, nil)
	//
	//} else if _network := asKCP(network); _network != "" {
	//	lis, err = kcp.InheritedListen(_network, laddr, tlsConfig, dataShards, parityShards)
	//
	//} else {

	lis, err = inherit.Listen(network, addrStr)
	if err == nil && tlsConfig != nil {
		if len(tlsConfig.Certificates) == 0 && tlsConfig.GetCertificate == nil {
			return nil, errors.New("tls: neither Certificates nor GetCertificate set in Config")
		}
		lis = tls.NewListener(lis, tlsConfig)
	}
	//}

	if err == nil {
		PushParentAddr(network, host, lis.Addr().String())
	}
	return
}
