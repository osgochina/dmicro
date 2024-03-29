package auth_test

import (
	"context"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/codec"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/plugin/auth"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := drpc.NewEndpoint(
			drpc.EndpointConfig{ListenPort: 9080},
			authChecker,
		)
		srv.RouteCall(new(Home))
		go srv.ListenAndServe()
		time.Sleep(1e9)
		cli := drpc.NewEndpoint(drpc.EndpointConfig{}, authBearer)
		sess, stat := cli.Dial("127.0.0.1:9080")
		if !stat.OK() {
			t.Fatal(stat)
		}
		var result interface{}
		stat = sess.Call("/home/test",
			map[string]string{
				"author": "clownfish",
			},
			&result,
			drpc.WithSetMeta("peer_id", "110"),
		).Status()
		if !stat.OK() {
			t.Error(stat)
		}
		t.Logf("result:%v", result)
		time.Sleep(3e9)
	})
}

func TestAuthBearer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := drpc.NewEndpoint(
			drpc.EndpointConfig{ListenPort: 9081},
			authChecker,
		)
		srv.RouteCall(new(Home))
		go srv.ListenAndServe()
		time.Sleep(10 * time.Minute)

	})
}

func TestAuthChecker(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cli := drpc.NewEndpoint(drpc.EndpointConfig{}, authBearer)
		sess, stat := cli.Dial("127.0.0.1:9080")
		if !stat.OK() {
			t.Fatal(stat)
		}
		var result interface{}
		stat = sess.Call("/home/test",
			map[string]string{
				"author": "clownfish",
			},
			&result,
			drpc.WithSetMeta("peer_id", "110"),
		).Status()
		if !stat.OK() {
			t.Error(stat)
		}
		t.Logf("result:%v", result)
		time.Sleep(3e9)
	})
}

const clientAuthInfo = "client-auth-info-12345"

var authBearer = auth.NewBearerPlugin(
	func(sess auth.Session, fn auth.SendOnce) *drpc.Status {
		var ret string
		stat := fn(clientAuthInfo, &ret)
		if !stat.OK() {
			return stat
		}
		internal.Infof(context.TODO(), "auth info: %s, result: %s", clientAuthInfo, ret)
		return nil
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

type Home struct {
	drpc.CallCtx
}

func (h *Home) Test(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
	return map[string]interface{}{
		"arg": *arg,
	}, nil
}
