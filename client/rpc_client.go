package client

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/message"
	"github.com/osgochina/dmicro/drpc/plugin/heartbeat"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/metrics"
	"github.com/osgochina/dmicro/registry"
	"github.com/osgochina/dmicro/selector"
	"sync"
	"time"
)

var (
	// DefaultBodyCodec 默认的消息编码
	defaultBodyCodec = "json"
	// DefaultSessionAge 默认session会话生命周期
	defaultSessionAge = time.Duration(0)
	// DefaultContextAge 默认单次请求生命周期
	defaultContextAge = time.Duration(0)
	// DefaultDialTimeout 作为客户端角色时，请求服务端的超时时间
	defaultDialTimeout = time.Second * 5
	// DefaultSlowCometDuration 慢处理定义时间
	defaultSlowCometDuration = time.Duration(0)
	// DefaultRetryTimes 默认重试次数
	defaultRetryTimes = 2

	// RerClientClosed 客户端已关闭错误信息
	rerClientClosed = drpc.NewStatus(100, "client is closed", "")
)

// RpcClient rpc客户端结构体
type RpcClient struct {
	endpoint drpc.Endpoint
	opts     Options
	closeCh  chan bool
	closeMu  sync.Mutex
}

// NewRpcClient 创建rpc客户端
func NewRpcClient(serviceName string, opt ...Option) *RpcClient {
	opts := NewOptions(append([]Option{OptServiceName(serviceName)}, opt...)...)
	//如果设置了心跳包，则发送心跳
	var heartbeatPing heartbeat.Ping
	if opts.HeartbeatTime > time.Duration(0) {
		heartbeatPing = heartbeat.NewPing(int(opts.HeartbeatTime/time.Second), false)
		opts.GlobalPlugin = append(opts.GlobalPlugin, heartbeatPing)
	}
	// 如果存在metrics组件，则获取该组件的rpc插件
	if opts.Metrics != nil {
		// 未设置 service name 才需要初始化
		if len(opts.Metrics.Options().ServiceName) == 0 {
			opts.Metrics.Init(metrics.OptServiceName(serviceName))
		}
		opts.GlobalPlugin = append(opts.GlobalPlugin, opts.Metrics.Options().Plugins...)
		opts.Metrics.Start()
	}
	endpoint := drpc.NewEndpoint(opts.EndpointConfig(), opts.GlobalPlugin...)
	// 优先使用已生成的证书对象
	if opts.TLSConfig != nil {
		endpoint.SetTLSConfig(opts.TLSConfig)
	} else {
		//如果设置了证书，则必须执行证书操作
		if len(opts.TlsCertFile) > 0 && len(opts.TlsKeyFile) > 0 {
			err := endpoint.SetTLSConfigFromFile(opts.TlsCertFile, opts.TlsKeyFile)
			if err != nil {
				logger.Fatalf(context.TODO(), "%v", err)
			}
		}
	}
	rc := &RpcClient{
		opts:     opts,
		endpoint: endpoint,
		closeCh:  make(chan bool),
	}
	return rc
}

// Options 获取配置信息
func (that *RpcClient) Options() Options {
	return that.opts
}

// Call 请求服务端
func (that *RpcClient) Call(serviceMethod string, args interface{}, result interface{}, setting ...message.MsgSetting) drpc.CallCmd {
	select {
	case <-that.closeCh:
		return drpc.NewFakeCallCmd(serviceMethod, args, result, rerClientClosed)
	default:
	}
	var (
		callCmd  drpc.CallCmd
		connFail bool
	)
	for i := 0; i < that.opts.RetryTimes; i++ {
		sess, stat := that.selectSession(serviceMethod)
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
			logger.Debugf(context.TODO(), "链接第[%d]出错，错误原因: %s", i, callCmd.Status().String())
		}
	}
	return callCmd
}

// Push 发送push消息
func (that *RpcClient) Push(serviceMethod string, arg interface{}, setting ...message.MsgSetting) *drpc.Status {
	select {
	case <-that.closeCh:
		return rerClientClosed
	default:
	}
	var (
		stat     *drpc.Status
		connFail bool
		sess     drpc.Session
	)
	for i := 0; i < that.opts.RetryTimes; i++ {
		sess, stat = that.selectSession(serviceMethod)
		if stat != nil {
			return stat
		}
		stat = sess.Push(serviceMethod, arg, setting...)
		connFail = !drpc.IsConnError(stat)
		if connFail {
			return stat
		}
		if i > 0 {
			logger.Debugf(context.TODO(), "链接第[%d]出错，错误原因: %s", i, stat.String())
		}
	}
	return stat
}

// AsyncCall 异步请求
func (that *RpcClient) AsyncCall(serviceMethod string, arg interface{}, result interface{}, callCmdChan chan<- drpc.CallCmd, setting ...message.MsgSetting) drpc.CallCmd {
	if callCmdChan == nil {
		callCmdChan = make(chan drpc.CallCmd, 10) // buffered.
	} else {
		if cap(callCmdChan) == 0 {
			logger.Panicf(context.TODO(), "*Client.AsyncCall(): callCmdChan channel is unbuffered")
		}
	}
	select {
	case <-that.closeCh:
		callCmd := drpc.NewFakeCallCmd(serviceMethod, arg, result, rerClientClosed)
		callCmdChan <- callCmd
		return callCmd
	default:
	}
	sess, stat := that.selectSession(serviceMethod)
	if stat != nil {
		callCmd := drpc.NewFakeCallCmd(serviceMethod, arg, result, stat)
		callCmdChan <- callCmd
		return callCmd
	}
	callCmd := sess.AsyncCall(serviceMethod, arg, result, callCmdChan, setting...)
	return callCmd
}

// SubRoute 设置服务的路由组
func (that *RpcClient) SubRoute(pathPrefix string, plugin ...drpc.Plugin) *drpc.SubRouter {
	return that.endpoint.SubRoute(pathPrefix, plugin...)
}

// RoutePush 使用结构体注册PUSH处理程序，并且返回地址
func (that *RpcClient) RoutePush(ctrlStruct interface{}, plugin ...drpc.Plugin) []string {
	return that.endpoint.RoutePush(ctrlStruct, plugin...)
}

// RoutePushFunc 使用函数注册PUSH处理程序，并且返回地址
func (that *RpcClient) RoutePushFunc(pushHandleFunc interface{}, plugin ...drpc.Plugin) string {
	return that.endpoint.RoutePushFunc(pushHandleFunc, plugin...)
}

// Endpoint 返回Endpoint对象
func (that *RpcClient) Endpoint() drpc.Endpoint {
	return that.endpoint
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
		if that.opts.Metrics != nil {
			that.opts.Metrics.Shutdown()
		}
	}
}

// 选择session
func (that *RpcClient) selectSession(serviceMethod string) (drpc.Session, *drpc.Status) {

	next, err := that.next(that.Options().ServiceName)
	if err != nil {
		return nil, err
	}
	node, e := next()
	if e != nil {
		if e == selector.ErrNotFound {
			return nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client service %s: %s", that.Options().ServiceName, e.Error()))
		}
		return nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client error selecting %s node: %s", that.Options().ServiceName, e.Error()))
	}
	addr := node.Address
	s, found := that.endpoint.GetSession(addr)
	if found && s.Health() {
		if s.Health() {
			return s, nil
		}
		_ = s.Close()
	}
	s, stat := that.endpoint.Dial(addr, that.opts.ProtoFunc)
	if !stat.OK() {
		return s, drpc.NewStatus(drpc.CodeDialFailed, "", stat)
	}
	s.SetID(addr)
	return s, nil
}

// 获取服务可用的节点列表
func (that *RpcClient) next(serviceName string) (selector.Next, *drpc.Status) {
	if that.opts.Selector == nil {
		that.defaultSelector(serviceName)
	}
	next, err := that.opts.Selector.Select(serviceName)
	if err != nil {
		if err == selector.ErrNotFound {
			return nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client service %s: %s", serviceName, err.Error()))
		}
		return nil, drpc.NewStatus(drpc.CodeInternalServerError, fmt.Sprintf("dmicro.client error selecting %s node: %s", serviceName, err.Error()))
	}

	return next, nil
}

// 如果没有设置selector,则初始化
func (that *RpcClient) defaultSelector(serviceName string) {
	that.opts.Registry = registry.DefaultRegistry
	err := that.opts.Registry.Init(registry.OptServiceName(serviceName))
	if err != nil {
		logger.Error(context.TODO(), err)
	}
	// 初始化默认selector
	if that.opts.Selector == nil {
		that.opts.Selector = selector.NewSelector(selector.OptRegistry(that.opts.Registry))
	} else {
		_ = that.opts.Selector.Init(selector.OptRegistry(that.opts.Registry))
	}
}
