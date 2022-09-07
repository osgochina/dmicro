package server

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/heartbeat"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/metric"
	"github.com/osgochina/dmicro/registry"
	"net"
	"sync"
	"time"
)

var (
	// defaultBodyCodec 默认的消息编码
	defaultBodyCodec = "json"
	// DefaultSessionAge 默认session会话生命周期
	defaultSessionAge = time.Duration(0)
	// DefaultContextAge 默认单次请求生命周期
	defaultContextAge = time.Duration(0)
	// DefaultSlowCometDuration 慢处理定义时间
	defaultSlowCometDuration = time.Duration(0)
)

// RpcServer rpc服务端
type RpcServer struct {
	endpoint drpc.Endpoint
	opts     Options
	closeCh  chan bool
	closeMu  sync.Mutex
}

// NewRpcServer 创建rpcServer
func NewRpcServer(serviceName string, opt ...Option) *RpcServer {

	opts := newOptions(append([]Option{OptServiceName(serviceName)}, opt...)...)
	//如果设置了心跳包，则发送心跳
	var heartbeatPong heartbeat.Pong
	if opts.EnableHeartbeat {
		heartbeatPong = heartbeat.NewPong()
		opts.GlobalPlugin = append(opts.GlobalPlugin, heartbeatPong)
	}
	// 增加服务注册中心组件
	reg := opts.Registry
	if reg == nil {
		reg = registry.DefaultRegistry
		// mdns组件需要初始化服务名称和版本
		_ = reg.Init(registry.ServiceName(opts.ServiceName), registry.ServiceVersion(opts.ServiceVersion))
	}
	// 如果存在metric组件，则获取该组件的rpc插件
	if opts.metric != nil {
		opts.metric.Init(metric.OptServiceName(serviceName))
		opts.GlobalPlugin = append(opts.GlobalPlugin, opts.metric.Options().Plugins...)
	}
	opts.GlobalPlugin = append(opts.GlobalPlugin, registry.NewRegistryPlugin(reg))
	endpoint := drpc.NewEndpoint(opts.EndpointConfig(), opts.GlobalPlugin...)
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
	rc := &RpcServer{
		opts:     opts,
		endpoint: endpoint,
		closeCh:  make(chan bool),
	}
	return rc
}

// Options 获取配置信息
func (that *RpcServer) Options() Options {
	return that.opts
}

// Endpoint 返回Endpoint对象
func (that *RpcServer) Endpoint() drpc.Endpoint {
	return that.endpoint
}

// Router 获取路由对象
func (that *RpcServer) Router() *drpc.Router {
	return that.endpoint.Router()
}

// SubRoute 添加分组
func (that *RpcServer) SubRoute(pathPrefix string, plugin ...drpc.Plugin) *drpc.SubRouter {
	return that.endpoint.SubRoute(pathPrefix, plugin...)
}

// RouteCall 通过Struct注册CALL处理方法,返回资源标识符
func (that *RpcServer) RouteCall(ctrlStruct interface{}, plugin ...drpc.Plugin) []string {
	return that.endpoint.RouteCall(ctrlStruct, plugin...)
}

// RouteCallFunc 通过func注册CALL处理方法，返回资源标识符
func (that *RpcServer) RouteCallFunc(callHandleFunc interface{}, plugin ...drpc.Plugin) string {
	return that.endpoint.RouteCallFunc(callHandleFunc, plugin...)
}

// RoutePush 通过Struct注册PUSH处理方法,返回资源标识符
func (that *RpcServer) RoutePush(ctrlStruct interface{}, plugin ...drpc.Plugin) []string {
	return that.endpoint.RoutePush(ctrlStruct, plugin...)
}

// RoutePushFunc 通过func注册PUSH处理方法，返回资源标识符
func (that *RpcServer) RoutePushFunc(pushHandleFunc interface{}, plugin ...drpc.Plugin) string {
	return that.endpoint.RoutePushFunc(pushHandleFunc, plugin...)
}

// SetUnknownCall 设置默认处理方法
// 当请求call类型的资源标识符不存在时，则执行其设置的方法
func (that *RpcServer) SetUnknownCall(fn func(drpc.UnknownCallCtx) (interface{}, *drpc.Status), plugin ...drpc.Plugin) {
	that.endpoint.SetUnknownCall(fn, plugin...)
}

// SetUnknownPush 设置默认处理方法
//
//	当请求push类型的资源标识符不存在时，则执行其设置的方法
func (that *RpcServer) SetUnknownPush(fn func(drpc.UnknownPushCtx) *drpc.Status, plugin ...drpc.Plugin) {
	that.endpoint.SetUnknownPush(fn, plugin...)
}

// Close 关闭客户端对象
func (that *RpcServer) Close() {
	that.closeMu.Lock()
	defer that.closeMu.Unlock()
	select {
	case <-that.closeCh:
		return
	default:
		close(that.closeCh)
		_ = that.endpoint.Close()
	}
}

// CountSession 返回回话数量
func (that *RpcServer) CountSession() int {
	return that.endpoint.CountSession()
}

// GetSession 通过session id获取session
func (that *RpcServer) GetSession(sessionId string) (drpc.Session, bool) {
	return that.endpoint.GetSession(sessionId)
}

// ListenAndServe 启动并监听服务
func (that *RpcServer) ListenAndServe(protoFunc ...proto.ProtoFunc) error {
	//如果未传入传输协议，则使用opts中配置的协议
	var protoFs []proto.ProtoFunc
	if len(protoFunc) > 0 {
		protoFs = protoFunc
	} else if that.opts.ProtoFunc != nil {
		protoFs = []proto.ProtoFunc{that.opts.ProtoFunc}
	}
	// 如果存在metric组件，则启动它
	if that.opts.metric != nil {
		that.opts.metric.Start()
	}
	return that.endpoint.ListenAndServe(protoFs...)
}

// ServeConn 传入指定的conn，生成session
// 提示：
// 1. 不支持断开链接后自动重拨
// 2. 不检查TLS
// 3. 执行 PostAcceptPlugin 插件
func (that *RpcServer) ServeConn(conn net.Conn, protoFunc ...proto.ProtoFunc) (drpc.Session, *drpc.Status) {
	return that.endpoint.ServeConn(conn, protoFunc...)
}
