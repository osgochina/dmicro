package websocket_test

import (
	"fmt"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/mixer/websocket"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/jsonSubProto"
	"github.com/osgochina/dmicro/logger"
	websocket2 "golang.org/x/net/websocket"
	"net/http"
	"testing"
	"time"
)

const clientAuthKey = "access_token"
const clientUserID = "user-1234"

var handshakePlugin = websocket.NewHandshakeAuthPlugin(
	func(r *http.Request) (sessionId string, status *drpc.Status) {
		token := websocket.QueryToken(clientAuthKey, r)
		logger.Infof("auth token: %v", token)
		if token != clientAuthInfo {
			return "", drpc.NewStatus(drpc.CodeUnauthorized, drpc.CodeText(drpc.CodeUnauthorized))
		}
		return clientUserID, nil
	},
	func(sess drpc.Session) *drpc.Status {
		logger.Infof("login userID: %v", sess.ID())
		return nil
	},
)

func TestHandshakeWebsocketAuth(t *testing.T) {
	srv := drpc.NewEndpoint(drpc.EndpointConfig{}, handshakePlugin)
	http.Handle("/token", websocket.NewJSONServeHandler(srv, nil))
	go http.ListenAndServe(":9094", nil)
	srv.RouteCall(new(P))
	time.Sleep(time.Millisecond * 200)

	// example in Browser: ws://localhost/token?access_token=clientAuthInfo
	rawQuery := fmt.Sprintf("/token?%s=%s", clientAuthKey, clientAuthInfo)
	cli := drpc.NewEndpoint(drpc.EndpointConfig{}, websocket.NewDialPlugin(rawQuery))
	sess, stat := cli.Dial(":9094", jsonSubProto.NewJSONSubProtoFunc())
	if !stat.OK() {
		t.Fatal(stat)
	}
	var result int
	stat = sess.Call("/p/divide", &Arg{
		A: 10,
		B: 2,
	}, &result,
	).Status()
	if !stat.OK() {
		t.Fatal(stat)
	}
	t.Logf("10/2=%d", result)
	time.Sleep(time.Millisecond * 200)

	// error test
	rawQuery = fmt.Sprintf("/token?%s=wrongToken", clientAuthKey)
	cli = drpc.NewEndpoint(drpc.EndpointConfig{}, websocket.NewDialPlugin(rawQuery))
	sess, stat = cli.Dial(":9094", jsonSubProto.NewJSONSubProtoFunc())
	if stat.OK() {
		t.Fatal("why dial correct by wrong token?")
	}
	time.Sleep(time.Millisecond * 200)
}

func TestHandshakeWebsocketAuthCustomizedHandshake(t *testing.T) {
	srv := websocket.NewServer("/token", drpc.EndpointConfig{ListenPort: 9095}, handshakePlugin)
	srv.RouteCall(new(P))
	srv.SetHandshake(func(config *websocket2.Config, request *http.Request) error {
		fmt.Println(config.Origin)
		fmt.Println(request.RequestURI)
		return nil
	})
	go srv.ListenAndServeJSON()
	time.Sleep(time.Millisecond * 200)

	// example in Browser: ws://localhost/token?access_token=clientAuthInfo
	rawQuery := fmt.Sprintf("/token?%s=%s", clientAuthKey, clientAuthInfo)
	cli := drpc.NewEndpoint(drpc.EndpointConfig{}, websocket.NewDialPlugin(rawQuery))
	sess, stat := cli.Dial(":9095", jsonSubProto.NewJSONSubProtoFunc())
	if !stat.OK() {
		t.Fatal(stat)
	}
	var result int
	stat = sess.Call("/p/divide", &Arg{
		A: 10,
		B: 2,
	}, &result,
	).Status()
	if !stat.OK() {
		t.Fatal(stat)
	}
	t.Logf("10/2=%d", result)
	time.Sleep(time.Millisecond * 200)

	// error test
	rawQuery = fmt.Sprintf("/token?%s=wrongToken", clientAuthKey)
	cli = drpc.NewEndpoint(drpc.EndpointConfig{}, websocket.NewDialPlugin(rawQuery))
	sess, stat = cli.Dial(":9095", jsonSubProto.NewJSONSubProtoFunc())
	if stat.OK() {
		t.Fatal("why dial correct by wrong token?")
	}
	time.Sleep(time.Millisecond * 200)
}
