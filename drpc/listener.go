package drpc

import (
	"crypto/tls"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/osgochina/dmicro/drpc/netproto/kcp"
	"github.com/osgochina/dmicro/drpc/netproto/normal"
	"github.com/osgochina/dmicro/drpc/netproto/quic"
	"github.com/osgochina/dmicro/utils"
	"net"
)

var testTLSConfig = utils.GenerateTLSConfigForServer()

// NewInheritedListener 创建一个支持优雅重启，支持继承监听的监听器
func NewInheritedListener(addr net.Addr, tlsConfig *tls.Config) (lis net.Listener, err error) {
	addrStr := addr.String()
	network := addr.Network()
	if _network := asQUIC(network); _network != "" {
		if tlsConfig == nil {
			tlsConfig = testTLSConfig
		}
		lis, err = quic.InheritedListen(_network, addrStr, tlsConfig, nil)

	} else if _network := asKCP(network); _network != "" {
		lis, err = kcp.InheritedListen(_network, addrStr, tlsConfig, kcp.DefaultDataShards, kcp.DefaultParityShards)
	} else {
		lis, err = normal.Listen(network, addrStr)
		if err == nil && tlsConfig != nil {
			if len(tlsConfig.Certificates) == 0 && tlsConfig.GetCertificate == nil {
				return nil, gerror.New("tls: neither Certificates nor GetCertificate set in Config")
			}
			lis = tls.NewListener(lis, tlsConfig)
		}
	}
	return
}
