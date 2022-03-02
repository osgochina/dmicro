package client

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/plugin/heartbeat"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/selector"
	"sync"
	"time"
)

var (
	// DefaultBodyCodec 默认的消息编码
	DefaultBodyCodec = "json"
	// DefaultSessionAge 默认session会话生命周期
	DefaultSessionAge = time.Duration(0)
	// DefaultContextAge 默认单次请求生命周期
	DefaultContextAge = time.Duration(0)
	// DefaultDialTimeout 作为客户端角色时，请求服务端的超时时间
	DefaultDialTimeout = time.Second * 5
	// DefaultSlowCometDuration 慢处理定义时间
	DefaultSlowCometDuration = time.Duration(0)
	// DefaultRetryTimes 默认重试次数
	DefaultRetryTimes = 2

	// RerClientClosed 客户端已关闭错误信息
	RerClientClosed = drpc.NewStatus(100, "client is closed", "")
)

type RpcClient struct {
	serviceName string // 服务名称
	endpoint    drpc.Endpoint
	opts        Options
	sessMap     *gmap.StrAnyMap
	closeCh     chan bool
	closeMu     sync.Mutex
}

func NewRpcClient(serviceName string, opt ...Option) *RpcClient {

	opts := NewOptions(opt...)
	//如果设置了心跳包，则发送心跳
	var heartbeatPing heartbeat.Ping
	if opts.HeartbeatTime > time.Duration(0) {
		heartbeatPing = heartbeat.NewPing(int(opts.HeartbeatTime/time.Second), false)
		opts.GlobalLeftPlugin = append(opts.GlobalLeftPlugin, heartbeatPing)
	}
	endpoint := drpc.NewEndpoint(opts.EndpointConfig(), opts.GlobalLeftPlugin...)
	// 优先使用已生成的证书对象
	if opts.TLSConfig != nil {
		endpoint.SetTLSConfig(opts.TLSConfig)
	} else {
		//如果设置了证书，则必须执行证书操作
		if len(opts.TlsCertFile) > 0 && len(opts.TlsKeyFile) > 0 {
			err := endpoint.SetTLSConfigFromFile(opts.TlsCertFile, opts.TlsKeyFile)
			if err != nil {
				logger.Fatalf("%v", err)
			}
		}
	}
	rc := &RpcClient{
		serviceName: serviceName,
		opts:        opts,
		endpoint:    endpoint,
		sessMap:     gmap.NewStrAnyMap(true),
	}
	return rc
}

func (that *RpcClient) Options() Options {
	return that.opts
}

// Call 请求服务端
func (that *RpcClient) Call(serviceMethod string, args interface{}, result interface{}, setting ...message.MsgSetting) drpc.CallCmd {
	select {
	case <-that.closeCh:
		return drpc.NewFakeCallCmd(serviceMethod, args, result, RerClientClosed)
	default:
	}
	var (
		callCmd  drpc.CallCmd
		connFail bool
	)
	for i := 0; i < that.opts.RetryTimes; i++ {
		addr, sess, stat := that.selectSession(serviceMethod)
		if stat != nil {
			return drpc.NewFakeCallCmd(serviceMethod, args, result, stat)
		}
		var callCmdChan = make(chan drpc.CallCmd, 1)
		sess.AsyncCall(serviceMethod, args, result, callCmdChan, setting...)
		callCmd = <-callCmdChan
		// 判断错误类型是否是链接出错，如果不是链接出错，则直接返回错误信息，如果是链接出错，则删除该链接，重新执行
		connFail = drpc.IsConnError(callCmd.Status())
		if !connFail {
			return callCmd
		}
		if i > 0 {
			logger.Debugf("链接第[%d]出错，错误原因: %s", i, callCmd.Status().String())
		}
		// 删除出错链接
		_ = sess.Close()
		that.sessMap.Remove(addr)
	}
	return callCmd
}

// Push 发送push消息
func (that *RpcClient) Push(serviceMethod string, arg interface{}, setting ...message.MsgSetting) *drpc.Status {
	select {
	case <-that.closeCh:
		return RerClientClosed
	default:
	}
	var (
		stat     *drpc.Status
		connFail bool
		sess     drpc.Session
		addr     string
	)
	for i := 0; i < that.opts.RetryTimes; i++ {
		addr, sess, stat = that.selectSession(serviceMethod)
		if stat != nil {
			return stat
		}
		stat = sess.Push(serviceMethod, arg, setting...)
		connFail = !drpc.IsConnError(stat)
		if connFail {
			return stat
		}
		if i > 0 {
			logger.Debugf("链接第[%d]出错，错误原因: %s", i, stat.String())
		}
		// 删除出错链接
		_ = sess.Close()
		that.sessMap.Remove(addr)
	}
	return stat
}

func (that *RpcClient) AsyncCall(serviceMethod string, arg interface{}, result interface{}, callCmdChan chan<- drpc.CallCmd, setting ...message.MsgSetting) drpc.CallCmd {
	if callCmdChan == nil {
		callCmdChan = make(chan drpc.CallCmd, 10) // buffered.
	} else {
		if cap(callCmdChan) == 0 {
			logger.Panicf("*Client.AsyncCall(): callCmdChan channel is unbuffered")
		}
	}
	select {
	case <-that.closeCh:
		callCmd := drpc.NewFakeCallCmd(serviceMethod, arg, result, RerClientClosed)
		callCmdChan <- callCmd
		return callCmd
	default:
	}

	addr, sess, stat := that.selectSession(serviceMethod)
	if stat != nil {
		callCmd := drpc.NewFakeCallCmd(serviceMethod, arg, result, stat)
		callCmdChan <- callCmd
		return callCmd
	}
	callCmd := sess.AsyncCall(serviceMethod, arg, result, callCmdChan, setting...)
	// 如果报链接错误，则删除缓存中的链接
	if drpc.IsConnError(callCmd.Status()) {
		_ = sess.Close()
		that.sessMap.Remove(addr)
	}
	return callCmd
}

// 选择session
func (that *RpcClient) selectSession(serviceMethod string) (string, drpc.Session, *drpc.Status) {

	next, err := that.next(that.serviceName)
	if err != nil {
		return "", nil, err
	}
	node, e := next()
	if e != nil {
		if e == selector.ErrNotFound {
			return "", nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client service %s: %s", that.serviceName, e.Error()))
		}
		return "", nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client error selecting %s node: %s", that.serviceName, e.Error()))
	}
	_s, found := that.sessMap.Search(node.Address)
	if !found {
		s, st := that.endpoint.Dial(node.Address, that.opts.ProtoFunc)
		if !st.OK() {
			return "", nil, st
		}
		that.sessMap.Set(node.Address, s)
		return node.Address, s, nil
	}

	return node.Address, _s.(drpc.Session), nil
}

// 获取服务可用的节点列表
func (that *RpcClient) next(serviceName string) (selector.Next, *drpc.Status) {
	next, err := that.opts.Selector.Select(serviceName)
	if err != nil {
		if err == selector.ErrNotFound {
			return nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client service %s: %s", serviceName, err.Error()))
		}
		return nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client error selecting %s node: %s", serviceName, err.Error()))
	}

	return next, nil
}

// Close 关闭客户端对象
func (that *RpcClient) Close() {
	that.closeMu.Lock()
	defer that.closeMu.Unlock()
	select {
	case <-that.closeCh:
		return
	default:
		close(that.closeCh)
		_ = that.endpoint.Close()
		_ = that.opts.Selector.Close()
		that.sessMap.Clear()
	}
}
