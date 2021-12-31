package kcp

import (
	"crypto/tls"
	"github.com/xtaci/kcp-go/v5"
	"net"
)

type (
	// UDPSession defines a KCP session implemented by UDP
	UDPSession = kcp.UDPSession
)

func DialAddrContext(network string, laddr *net.UDPAddr, raddr string, tlsConf *tls.Config, dataShards, parityShards int) (net.Conn, error) {
	host, port, err := net.SplitHostPort(raddr)
	if err != nil {
		return nil, err
	}
	if host == "" {
		raddr = "127.0.0.1:" + port
	}
	addr, err := net.ResolveUDPAddr(network, raddr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP(network, laddr)
	if err != nil {
		return nil, err
	}
	conn, err := kcp.NewConn2(addr, nil, dataShards, parityShards, udpConn)
	if err != nil {
		return nil, err
	}
	if tlsConf != nil {
		return tls.Client(conn, tlsConf), nil
	}
	return conn, nil
}

type Listener struct {
	*kcp.Listener
	tlsConf *tls.Config
	conn    net.PacketConn
}

var _ net.Listener = new(Listener)

func (that *Listener) Accept() (net.Conn, error) {
	conn, err := that.Listener.Accept()
	if err != nil {
		return nil, err
	}
	if that.tlsConf == nil {
		return conn, nil
	}
	return tls.Server(conn, that.tlsConf), nil
}

func (that *Listener) PacketConn() net.PacketConn {
	return that.conn
}

func ListenAddr(network, addr string, tlsConf *tls.Config, dataShards, parityShards int) (*Listener, error) {
	udpAddr, err := net.ResolveUDPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	return ListenUDPAddr(network, udpAddr, tlsConf, dataShards, parityShards)
}

func ListenUDPAddr(network string, udpAddr *net.UDPAddr, tlsConf *tls.Config, dataShards, parityShards int) (*Listener, error) {
	var conn net.PacketConn
	conn, err := net.ListenUDP(network, udpAddr)
	if err != nil {
		return nil, err
	}
	return Listen(conn, tlsConf, dataShards, parityShards)
}

func Listen(conn net.PacketConn, tlsConf *tls.Config, dataShards, parityShards int) (*Listener, error) {
	lis, err := kcp.ServeConn(nil, dataShards, parityShards, conn)
	if err != nil {
		return nil, err
	}
	return &Listener{Listener: lis, tlsConf: tlsConf, conn: conn}, nil
}
