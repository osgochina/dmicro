package websocket_test

import (
	"context"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/mixer/websocket"
	"github.com/osgochina/dmicro/drpc/plugin/auth"
	"testing"
	"time"
)

func TestJSONWebsocketAuth(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := websocket.NewServer(
			"/auth",
			drpc.EndpointConfig{ListenPort: 9093},
			authChecker,
		)
		srv.RouteCall(new(P))
		go srv.ListenAndServe()

		time.Sleep(time.Second * 1)

		cli := websocket.NewClient(
			"/auth",
			drpc.EndpointConfig{},
			authBearer,
		)
		sess, stat := cli.Dial(":9093")
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
		time.Sleep(time.Second)
	})

}

const clientAuthInfo = "client-auth-info-12345"

var authBearer = auth.NewBearerPlugin(
	func(sess auth.Session, fn auth.SendOnce) (stat *drpc.Status) {
		var ret string
		stat = fn(clientAuthInfo, &ret)
		if !stat.OK() {
			return
		}
		internal.Infof(context.TODO(), "auth info: %s, result: %s", clientAuthInfo, ret)
		return
	},
	drpc.WithBodyCodec(codec.PlainName),
)

var authChecker = auth.NewCheckerPlugin(
	func(sess auth.Session, fn auth.ReCvOnce) (ret interface{}, stat *drpc.Status) {
		var authInfo string
		stat = fn(&authInfo)
		if !stat.OK() {
			return
		}
		internal.Infof(context.TODO(), "auth info: %v", authInfo)
		if clientAuthInfo != authInfo {
			return nil, drpc.NewStatus(403, "auth fail", "auth fail detail")
		}
		return "pass", nil
	},
	drpc.WithBodyCodec(codec.PlainName),
)
