package websocket

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto"
	"golang.org/x/net/websocket"
	"net"
)

// NewDialPlugin 创建链接插件
func NewDialPlugin(rootPath string) drpc.Plugin {
	return &clientPlugin{fixRootPath(rootPath)}
}

// 链接到服务端后需要执行该插件
type clientPlugin struct {
	rootPath string
}

var (
	_ drpc.AfterDialPlugin = new(clientPlugin)
)

// Name 插件名
func (*clientPlugin) Name() string {
	return "websocket"
}

// AfterDial 不同于tcp，链接成功后就能直接发送数据，websocket的链接是从 tcp=>http=>websocket，所以链接成功后需要做一些逻辑处理
func (that *clientPlugin) AfterDial(sess drpc.EarlySession, isRedial bool) (stat *drpc.Status) {
	var location, origin string
	if sess.Endpoint().TLSConfig() == nil {
		location = fmt.Sprintf("ws://%s%s", sess.RemoteAddr().String(), that.rootPath)
		origin = fmt.Sprintf("ws://%s%s", sess.LocalAddr().String(), that.rootPath)
	} else {
		location = fmt.Sprintf("wss://%s%s", sess.RemoteAddr().String(), that.rootPath)
		origin = fmt.Sprintf("wss://%s%s", sess.LocalAddr().String(), that.rootPath)
	}
	// 生成websocket的配置信息
	cfg, err := websocket.NewConfig(location, origin)
	if err != nil {
		return drpc.NewStatus(drpc.CodeDialFailed, "upgrade to websocket failed", err.Error())
	}
	// 针对session底层使用的socket链接进程修改，把该链接对象上的websocket协议配置添加上去，并且把协议替换成新的wsProto
	sess.ModifySocket(func(conn net.Conn) (modifiedConn net.Conn, newProtoFunc proto.ProtoFunc) {
		conn, err := websocket.NewClient(cfg, conn)
		if err != nil {
			stat = drpc.NewStatus(drpc.CodeDialFailed, "upgrade to websocket failed", err.Error())
			return nil, nil
		}
		//如果是重连，说明协议已经添加了，不在需要再次添加
		if isRedial {
			return conn, sess.GetProtoFunc()
		}
		// 在设定的json协议或pb协议外层，在包一层wsProto协议
		return conn, NewWsProtoFunc(sess.GetProtoFunc())
	})
	return stat
}
