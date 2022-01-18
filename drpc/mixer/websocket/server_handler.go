package websocket

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/jsonSubProto"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/pbSubProto"
	"github.com/osgochina/dmicro/drpc/proto"
	"golang.org/x/net/websocket"
	"net/http"
	"net/url"
)

// 服务处理器对象
type serverHandler struct {
	endpoint  drpc.Endpoint
	protoFunc proto.ProtoFunc
	*websocket.Server
}

// NewJSONServeHandler 创建使用json格式编解码的处理器
func NewJSONServeHandler(endpoint drpc.Endpoint, handshake func(*websocket.Config, *http.Request) error) http.Handler {
	return NewServeHandler(endpoint, handshake, jsonSubProto.NewJSONSubProtoFunc())
}

// NewPbServeHandler 创建使用protobuf协议编解码的处理器
func NewPbServeHandler(endpoint drpc.Endpoint, handshake func(*websocket.Config, *http.Request) error) http.Handler {
	return NewServeHandler(endpoint, handshake, pbSubProto.NewPbSubProtoFunc())
}

// NewServeHandler 创建处理器
func NewServeHandler(endpoint drpc.Endpoint, handshake func(*websocket.Config, *http.Request) error, protoFunc ...proto.ProtoFunc) http.Handler {
	w := &serverHandler{
		endpoint:  endpoint,
		Server:    new(websocket.Server),
		protoFunc: NewWsProtoFunc(protoFunc...),
	}
	var scheme string
	if endpoint.TLSConfig() == nil {
		scheme = "ws"
	} else {
		scheme = "wss"
	}
	//创建websocket握手处理器
	w.Server.Handshake = func(cfg *websocket.Config, r *http.Request) error {
		cfg.Origin = &url.URL{
			Host:   r.RemoteAddr,
			Scheme: scheme,
		}
		// 先执行握手钩子
		if stat := w.beforeHandshake(r); !stat.OK() {
			return stat.Cause()
		}
		// 在执行自定义握手函数
		if handshake != nil {
			return handshake(cfg, r)
		}
		return nil
	}
	// 绑定处理器
	w.Server.Handler = w.handler
	// 支持tls
	w.Server.Config = websocket.Config{
		TlsConfig: endpoint.TLSConfig(),
	}
	return w
}

// websocket握手成功后，执行该方法
func (that *serverHandler) handler(conn *websocket.Conn) {
	// 把websocket的底层tcp链接传给drpc框架，方便它读取数据
	sess, err := that.endpoint.ServeConn(conn, that.protoFunc)
	if err != nil {
		internal.Errorf("serverHandler: %v", err)
		return
	}
	// 执行钩子
	if stat := that.afterAccept(sess, conn); !stat.OK() {
		if err := sess.Close(); err != nil {
			internal.Errorf("sess.Close(): %v", err)
		}
		return
	}
	// 等待会话关闭，结束该协程
	<-sess.CloseNotify()
}

// websocket握手之前需要执行该方法，如果处理后返回错误，则握手不成功，返回的错误为空则握手成功
func (that *serverHandler) beforeHandshake(r *http.Request) (stat *drpc.Status) {
	var pluginName string
	p := that.endpoint.PluginContainer()
	defer func() {
		if p := recover(); p != nil {
			internal.Errorf("[BeforeWebsocketHandshakePlugin:%s] addr:%s, panic:%v", pluginName, r.RemoteAddr, p)
			stat = statInternalServerError.Copy(p)
		}
	}()
	// 执行握手钩子
	for _, plugin := range p.GetAll() {
		if _plugin, ok := plugin.(BeforeWebsocketHandshakePlugin); ok {
			pluginName = plugin.Name()
			if stat = _plugin.BeforeHandshake(r); !stat.OK() {
				internal.Debugf("[BeforeWebsocketHandshakePlugin:%s] addr:%s, error:%s", pluginName, r.RemoteAddr, stat.String())
				return stat
			}
		}
	}
	return nil
}

// websocket握手成功后回调该方法，表示接收了该链接
func (that *serverHandler) afterAccept(sess drpc.Session, conn *websocket.Conn) (stat *drpc.Status) {
	var pluginName string
	p := that.endpoint.PluginContainer()
	defer func() {
		if p := recover(); p != nil {
			internal.Errorf("[AfterWebsocketAcceptPlugin:%s] network:%s, addr:%s, panic:%v", pluginName, sess.RemoteAddr().Network(), sess.RemoteAddr().String(), p)
			stat = statInternalServerError.Copy(p)
		}
	}()
	// 查找是否有定义的钩子，有的话执行该钩子
	for _, plugin := range p.GetAll() {
		if _plugin, ok := plugin.(AfterWebsocketAcceptPlugin); ok {
			pluginName = plugin.Name()
			if stat = _plugin.AfterAccept(sess, conn); !stat.OK() {
				internal.Debugf("[AfterWebsocketAcceptPlugin:%s] network:%s, addr:%s, error:%s", pluginName, sess.RemoteAddr().Network(), sess.RemoteAddr().String(), stat.String())
				return stat
			}
		}
	}
	return nil
}
