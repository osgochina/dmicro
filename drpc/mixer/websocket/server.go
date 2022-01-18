package websocket

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/jsonSubProto"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/pbSubProto"
	"github.com/osgochina/dmicro/drpc/proto"
	"github.com/osgochina/dmicro/utils"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
)

// 定义服务器内部错误
var statInternalServerError = drpc.NewStatus(drpc.CodeInternalServerError, drpc.CodeText(drpc.CodeInternalServerError), "")

// 定义服务插件
type (
	// BeforeWebsocketHandshakePlugin 在websocket握手之前执行该插件
	BeforeWebsocketHandshakePlugin interface {
		drpc.Plugin
		BeforeHandshake(r *http.Request) *drpc.Status
	}
	// AfterWebsocketAcceptPlugin 在接收到websocket链接之后执行该插件
	AfterWebsocketAcceptPlugin interface {
		drpc.Plugin
		AfterAccept(sess drpc.Session, conn *websocket.Conn) *drpc.Status
	}
)

// Server websocket 服务对象
type Server struct {
	drpc.Endpoint
	cfg       drpc.EndpointConfig
	rootPath  string                                       //要绑定的url地址
	lis       net.Listener                                 //监听的原始socket
	lisAddr   net.Addr                                     //绑定的地址
	serveMux  *http.ServeMux                               // http路由表结构
	server    *http.Server                                 //http对象
	handshake func(*websocket.Config, *http.Request) error // 自定义的握手回调方法
}

// NewServer 创建服务对象
// rootPath: 要绑定的url地址
// cfg: drpc的配置信息
// globalLeftPlugin: drpc的支持的插件
func NewServer(rootPath string, cfg drpc.EndpointConfig, globalLeftPlugin ...drpc.Plugin) *Server {
	e := drpc.NewEndpoint(cfg, globalLeftPlugin...)
	serveMux := http.NewServeMux()
	lisAddr := cfg.ListenAddr()
	host, port, _ := net.SplitHostPort(lisAddr.String())
	if port == "0" {
		if e.TLSConfig() != nil {
			port = "https"
		} else {
			port = "http"
		}
		lisAddr = utils.NewFakeAddr(lisAddr.Network(), host, port)
	}
	return &Server{
		Endpoint: e,
		cfg:      cfg,
		serveMux: serveMux,
		rootPath: fixRootPath(rootPath),
		lisAddr:  lisAddr,
		server:   &http.Server{Addr: lisAddr.String(), Handler: serveMux},
	}
}

// ListenAndServeJSON 使用json协议传输数据
func (that *Server) ListenAndServeJSON() error {
	return that.ListenAndServe(jsonSubProto.NewJSONSubProtoFunc())
}

// ListenAndServeProtobuf 使用protobuf协议传输数据
func (that *Server) ListenAndServeProtobuf() error {
	return that.ListenAndServe(pbSubProto.NewPbSubProtoFunc())
}

// ListenAndServe 监听地址端口并提供服务
func (that *Server) ListenAndServe(protoFunc ...proto.ProtoFunc) (err error) {
	network := that.cfg.Network
	switch network {
	default:
		return gerror.New("invalid network config, refer to the following: tcp, tcp4, tcp6")
	case "tcp", "tcp4", "tcp6":
	}
	// 绑定路由
	that.Handle(that.rootPath, NewServeHandler(that.Endpoint, that.handshake, protoFunc...))
	// 监听
	that.lis, err = drpc.NewInheritedListener(that.lisAddr, that.Endpoint.TLSConfig())
	if err != nil {
		return err
	}
	that.lisAddr = that.lis.Addr()
	internal.Printf("listen and serve (network:%s, addr:%s)", network, that.lisAddr)

	// 执行listen钩子
	for _, v := range that.Endpoint.PluginContainer().GetAll() {
		if e, ok := v.(drpc.AfterListenPlugin); ok {
			_ = e.AfterListen(that.lis.Addr())
		}
	}
	// 把tpc协议的监听地址传给http服务
	return that.server.Serve(that.lis)
}

// SetHandshake 设置自定义的握手方法
func (that *Server) SetHandshake(handshake func(*websocket.Config, *http.Request) error) {
	that.handshake = handshake
}

// Handle 绑定路由handler
func (that *Server) Handle(rootPath string, handler http.Handler) {
	that.serveMux.Handle(rootPath, handler)
}

// HandleFunc 绑定路由到匿名方法
func (that *Server) HandleFunc(rootPath string, handler func(http.ResponseWriter, *http.Request)) {
	that.serveMux.HandleFunc(rootPath, handler)
}

// Close 关闭服务
func (that *Server) Close() error {
	// 先关闭http服务，在关闭tcp层的endpoint
	err := that.server.Shutdown(context.Background())
	if err != nil {
		_ = that.Endpoint.Close()
		return err
	}
	return that.Endpoint.Close()
}
