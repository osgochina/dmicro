package quic

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"net"
)

type Listener struct {
	lis  quic.Listener
	conn net.PacketConn
}

var _ net.Listener = new(Listener)

// DialAddrContext 使用quic协议链接远端
// ctx: 上下文
// network: 网络类型,可选："udp", "udp4", "udp6"
// laddr: 本地监听的UDP协议地址
// raddr: 远端的地址
// tlsConf: 必须传入证书信息
// config：quic的配置信息,可以为nil
func DialAddrContext(ctx context.Context, network string, laddr *net.UDPAddr, raddr string, tlsConf *tls.Config, config *quic.Config) (net.Conn, error) {
	host, port, err := net.SplitHostPort(raddr)
	if err != nil {
		return nil, err
	}
	if host == "" {
		raddr = "127.0.0.1:" + port
	}
	udpAddr, err := net.ResolveUDPAddr(network, raddr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP(network, laddr)
	if err != nil {
		return nil, err
	}
	sess, err := quic.DialContext(ctx, udpConn, udpAddr, raddr, tlsConf, config)
	if err != nil {
		return nil, err
	}
	stream, err := sess.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{
		sess:   sess,
		stream: stream,
	}, nil
}

// ListenAddr 监听指定地址
// network: 网络类型,可选："udp", "udp4", "udp6"
// addr: 地址
// tlsConf: 必须传入证书信息
// config：quic的配置信息,可以为nil
func ListenAddr(network, addr string, tlsConf *tls.Config, config *quic.Config) (*Listener, error) {
	udpAddr, err := net.ResolveUDPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	return ListenUDPAddr(network, udpAddr, tlsConf, config)
}

// ListenUDPAddr 监听UDP协议地址
// network: 网络类型,可选："udp", "udp4", "udp6"
// addr: UDP协议地址
// tlsConf: 必须传入证书信息
// config：quic的配置信息,可以为nil
func ListenUDPAddr(network string, udpAddr *net.UDPAddr, tlsConf *tls.Config, config *quic.Config) (*Listener, error) {
	conn, err := net.ListenUDP(network, udpAddr)
	if err != nil {
		return nil, err
	}
	return Listen(conn, tlsConf, config)
}

// Listen 监听指定的链接
// conn: PacketConn类型的链接
// tlsConf: 必须传入证书信息
// config：quic的配置信息,可以为nil
func Listen(conn net.PacketConn, tlsConf *tls.Config, config *quic.Config) (*Listener, error) {
	if config == nil {
		config = &quic.Config{KeepAlive: true}
	}
	lis, err := quic.Listen(conn, tlsConf, config)
	if err != nil {
		return nil, err
	}
	return &Listener{
		lis:  lis,
		conn: conn,
	}, nil
}

func (that *Listener) PacketConn() net.PacketConn {
	return that.conn
}

func (that *Listener) Accept() (net.Conn, error) {
	ctx := context.TODO()
	sess, err := that.lis.Accept(ctx)
	if err != nil {
		return nil, err
	}
	stream, err := sess.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{
		sess:   sess,
		stream: stream,
	}, nil
}

func (that *Listener) Close() error {
	return that.lis.Close()
}

func (that *Listener) Addr() net.Addr {
	return that.lis.Addr()
}
