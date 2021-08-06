package drpc

import (
	"errors"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/socket"
	"math"
	"net"
	"strconv"
	"time"
)

// EndpointConfig 端点的配置
type EndpointConfig struct {

	// 网络类型; tcp, tcp4, tcp6, unix, unixpacket, kcp or quic"
	Network string

	//作为服务器角色时候，本地监听地址
	listenAddr net.Addr

	//当前要监听的服务器本地IP
	LocalIP string

	//需要监听的端口号
	ListenPort uint16

	// 默认的消息体编码格式
	DefaultBodyCodec string
	//默认会话生命周期
	DefaultSessionAge time.Duration
	//默认单次请求生命周期
	DefaultContextAge time.Duration
	//外部配置慢处理定义时间
	SlowCometDuration time.Duration
	//慢处理定义时间
	slowCometDuration time.Duration

	//是否打印会话中请求的 body或 metadata
	PrintDetail bool
	// 是否统计消耗时间
	CountTime bool

	// 作为客户端角色时，请求服务端的超时时间
	DialTimeout time.Duration
	//作为客户端角色时，请求服务端时候，本地使用的地址端口
	localAddr net.Addr
	// 在链接中断时候，试图链接服务端的最大重试次数。仅限客户端角色使用
	RedialTimes int
	//仅限客户端角色使用 试图链接服务端时候，重试的时间间隔
	RedialInterval time.Duration

	checked bool
}

// ListenAddr 获取本地监听地址，服务器角色
func (that *EndpointConfig) ListenAddr() net.Addr {
	_ = that.check()
	return that.listenAddr
}

// LocalAddr 获取本地地址，客户端角色
func (that *EndpointConfig) LocalAddr() net.Addr {
	_ = that.check()
	return that.localAddr
}

func (that *EndpointConfig) check() (err error) {
	if that.checked {
		return nil
	}
	that.checked = true

	if that.Network == "" {
		that.Network = "tcp"
	}

	if that.LocalIP == "" {
		that.LocalIP = "0.0.0.0"
	}

	//先初始化一个基础的本地地址
	that.localAddr, err = that.newAddr("0")
	if err != nil {
		return err
	}
	listenPort := strconv.FormatUint(uint64(that.ListenPort), 10)

	//获取本地监听地址，
	//network 可以赋值也可以使用默认
	//LocalIP 可以赋值也可以使用默认
	//listenPort 可以使用ListenPort赋值得到的，也可以使用0，随机获取
	that.listenAddr = NewFakeAddr(that.Network, that.LocalIP, listenPort)

	//慢请求的配置值，请求消耗时间大于该值，被定义为慢请求，默认为最大数字
	that.slowCometDuration = math.MaxInt64
	//如果外部设置了该事件，则使用外部设置时间
	if that.SlowCometDuration > 0 {
		that.slowCometDuration = that.SlowCometDuration
	}

	if len(that.DefaultBodyCodec) == 0 {
		that.DefaultBodyCodec = DefaultBodyCodec().Name()
	}

	//作为客户端，链接服务器重试间隔时间，默认为100毫秒
	if that.RedialInterval <= 0 {
		that.RedialInterval = time.Millisecond * 100
	}

	return nil
}

//返回网络地址
func (that *EndpointConfig) newAddr(port string) (net.Addr, error) {
	switch that.Network {
	default:
		return nil, errors.New(" Invalid network config, refer to the following: tcp, tcp4, tcp6, unix, unixpacket, kcp or quic")
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(that.Network, net.JoinHostPort(that.LocalIP, port))
	case "unix", "unixpacket":
		return net.ResolveUnixAddr(that.Network, net.JoinHostPort(that.LocalIP, port))
	case "kcp", "udp", "udp4", "udp6", "quic":
		udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(that.LocalIP, port))
		if err != nil {
			return nil, err
		}
		network := "kcp"
		if that.Network == "quic" {
			network = "quic"
		}
		return &FakeAddr{
			network: network,
			addr:    udpAddr.String(),
			host:    that.LocalIP,
			port:    strconv.Itoa(udpAddr.Port),
			udpAddr: udpAddr,
		}, nil
	}
}

func asQUIC(network string) string {
	switch network {
	case "quic":
		return "udp"
	default:
		return ""
	}
}

func asKCP(network string) string {
	switch network {
	case "kcp":
		return "udp"
	case "udp", "udp4", "udp6":
		return network
	default:
		return ""
	}
}

//默认消息体编码格式
var defaultBodyCodec codec.Codec = new(codec.JSONCodec)

func DefaultBodyCodec() codec.Codec {
	return defaultBodyCodec
}

// SetDefaultBodyCodec 设置默认消息体编码格式
func SetDefaultBodyCodec(codecID byte) error {
	c, err := codec.Get(codecID)
	if err != nil {
		return err
	}
	defaultBodyCodec = c
	return nil
}

// DefaultProtoFunc 默认传输协议
var DefaultProtoFunc = socket.DefaultProtoFunc

// SetDefaultProtoFunc 设置默认传输协议
var SetDefaultProtoFunc = socket.SetDefaultProtoFunc

// GetReadLimit 获取消息最大长度限制
var GetReadLimit = message.MsgSizeLimit

// SetReadLimit 设置消息最大长度限制
var SetReadLimit = message.SetMsgSizeLimit

// SetSocketKeepAlive 开启关闭链接保活
var SetSocketKeepAlive = socket.SetKeepAlive

// SetSocketKeepAlivePeriod 链接保活间隔时间
var SetSocketKeepAlivePeriod = socket.SetKeepAlivePeriod

// SocketReadBuffer 获取链接读缓冲区长度
var SocketReadBuffer = socket.ReadBuffer

// SetSocketReadBuffer 设置链接读缓冲区长度
var SetSocketReadBuffer = socket.SetReadBuffer

// SocketWriteBuffer 获取链接写缓冲区长度
var SocketWriteBuffer = socket.WriteBuffer

// SetSocketWriteBuffer 设置链接写缓冲区长度
var SetSocketWriteBuffer = socket.SetWriteBuffer

// SetSocketNoDelay 开启关闭no delay算法
var SetSocketNoDelay = socket.SetNoDelay
