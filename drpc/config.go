package drpc

import (
	"errors"
	"github.com/gogf/gf/util/gconv"
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
	Network string `json:"network" comment:"网络类型; tcp, tcp4, tcp6, unix, unixpacket"`
	//作为服务端角色时，要监听的服务器本地IP
	ListenIP string `json:"listen_ip" comment:"作为服务端角色时，要监听的服务器本地IP"`
	//作为服务端角色时，需要监听的本地端口号
	ListenPort uint16 `json:"listen_port" comment:"作为服务端角色时，需要监听的本地端口号"`
	//作为服务器角色时候，本地监听地址
	listenAddr net.Addr

	//作为客户端角色时,请求服务端时候，本地使用的地址
	LocalIP string `json:"local_ip" comment:"作为客户端角色时,请求服务端时候，本地使用的地址"`
	//作为客户端角色时,请求服务端时候，本地使用的地址端口号
	LocalPort uint16 `json:"local_port" comment:"作为客户端角色时,请求服务端时候，本地使用的地址端口号"`
	//作为客户端角色时，请求服务端时候，本地使用的地址端口
	localAddr net.Addr

	// 默认的消息体编码格式
	DefaultBodyCodec string `json:"default_body_codec" comment:"默认的消息体编码格式"`
	//默认session会话生命周期
	DefaultSessionAge time.Duration `json:"default_session_age" comment:"默认session会话生命周期"`
	//默认单次请求生命周期
	DefaultContextAge time.Duration `json:"default_context_age" comment:"默认单次请求生命周期"`
	//外部配置慢处理定义时间
	SlowCometDuration time.Duration `json:"slow_comet_duration" comment:"慢处理定义时长"`
	//慢处理定义时间
	slowCometDuration time.Duration

	//是否打印会话中请求的 body或 metadata
	PrintDetail bool `json:"print_detail" comment:"是否打印请求的详细信息，body和metadata"`
	// 是否统计消耗时间
	CountTime bool `json:"count_time" comment:"是否统计请求消耗时间"`

	// 作为客户端角色时，请求服务端的超时时间
	DialTimeout time.Duration `json:"dial_timeout" comment:"作为客户端角色时，请求服务端的超时时间"`
	// 仅限客户端角色使用,链接中断时候，试图链接服务端的最大重试次数。
	RedialTimes int `json:"redial_times" comment:"仅限客户端角色使用,链接中断时候，试图链接服务端的最大重试次数。"`
	//仅限客户端角色使用 试图链接服务端时候，重试的时间间隔.
	RedialInterval time.Duration `json:"redial_interval" comment:"仅限客户端角色使用 试图链接服务端时候，重试的时间间隔."`

	//该配置是否已经初始化检查
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

// 初始化检查，判断配置是否有误
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

	if that.LocalPort <= 0 {
		that.LocalPort = 0
	}

	//如果监听的ip为空，且设置的LocalIP不为空，则默认使用LocalIP，这样就可以兼容低版本只有LocalIP的情况
	if that.ListenIP == "" && that.LocalIP != "" {
		that.ListenIP = that.LocalIP
	}

	if that.ListenPort <= 0 {
		that.ListenPort = 0
	}

	//先初始化一个基础的本地地址，LocalPort为0的时候表示随机获取
	that.localAddr, err = that.newAddr(that.Network, that.LocalIP, gconv.String(that.LocalPort))
	if err != nil {
		return err
	}
	//获取本地监听地址，
	//network 可以赋值也可以使用默认
	//LocalIP 可以赋值也可以使用默认
	//ListenPort 可以使用0，随机获取
	that.listenAddr = NewFakeAddr(that.Network, that.ListenIP, gconv.String(that.ListenPort))

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
func (that *EndpointConfig) newAddr(network, ip, port string) (net.Addr, error) {
	switch network {
	default:
		return nil, errors.New(" Invalid network config, refer to the following: tcp, tcp4, tcp6, unix, unixpacket, kcp or quic")
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, net.JoinHostPort(ip, port))
	case "unix", "unixpacket":
		return net.ResolveUnixAddr(network, net.JoinHostPort(ip, port))
	case "kcp", "udp", "udp4", "udp6", "quic":
		udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(ip, port))
		if err != nil {
			return nil, err
		}
		n := "kcp"
		if network == "quic" {
			n = "quic"
		}
		return &FakeAddr{
			network: n,
			addr:    udpAddr.String(),
			host:    ip,
			port:    strconv.Itoa(udpAddr.Port),
			udpAddr: udpAddr,
		}, nil
	}
}

//默认消息体编码格式
var defaultBodyCodec codec.Codec = new(codec.JSONCodec)

// DefaultBodyCodec 获取当前默认消息体编码格式
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

// SetSocketKeepAlive 开启关闭死链检测
var SetSocketKeepAlive = socket.SetKeepAlive

// SetSocketKeepAlivePeriod 死链检测间隔时间
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
