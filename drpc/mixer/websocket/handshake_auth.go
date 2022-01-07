package websocket

import (
	"github.com/osgochina/dmicro/drpc"
	"golang.org/x/net/websocket"
	"net/http"
)

// NewHandshakeAuthPlugin 创建握手权限检查插件
func NewHandshakeAuthPlugin(ckFn Checker, apFn Acceptor) *HandshakeAuthPlugin {
	return &HandshakeAuthPlugin{
		CheckFunc:  ckFn,
		AcceptFunc: apFn,
	}
}

// Checker 定义检查方法，需要业务自己实现它,检查成功后会返回会话id，该id可以自定义
type Checker func(r *http.Request) (sessionID string, status *drpc.Status)

// Acceptor 检查通过后，接收该链接后，会调用该方法，需要业务自己实现它，传入的参数是会话对象，可以在此做一些初始化操作，
type Acceptor func(sess drpc.Session) *drpc.Status

// HandshakeAuthPlugin 握手权限检查插件
type HandshakeAuthPlugin struct {
	CheckFunc  Checker
	AcceptFunc Acceptor
}

var (
	_ AfterWebsocketAcceptPlugin     = new(HandshakeAuthPlugin)
	_ BeforeWebsocketHandshakePlugin = new(HandshakeAuthPlugin)
)

func (that *HandshakeAuthPlugin) Name() string {
	return "handshake-auth-plugin"
}

const sessionHeader = "Drpc-Session-Id"

// BeforeHandshake 握手之前回调该方法，可以在此做权限认证的操作
func (that *HandshakeAuthPlugin) BeforeHandshake(r *http.Request) *drpc.Status {
	if that.CheckFunc == nil {
		return nil
	}
	id, stat := that.CheckFunc(r)
	r.Header.Set(sessionHeader, id)
	return stat
}

// AfterAccept 握手成功后，接受该链接后回调此方法
func (that *HandshakeAuthPlugin) AfterAccept(sess drpc.Session, conn *websocket.Conn) *drpc.Status {
	if that.AcceptFunc == nil {
		return nil
	}
	id := conn.Request().Header.Get(sessionHeader)
	sess.SetID(id)
	stat := that.AcceptFunc(sess)
	return stat
}

// QueryToken 工具函数，方便获取认证用的token
func QueryToken(tokenKey string, r *http.Request) (token string) {
	queryParams := r.URL.Query()
	if values, ok := queryParams[tokenKey]; ok {
		token = values[0]
	}
	return token
}
