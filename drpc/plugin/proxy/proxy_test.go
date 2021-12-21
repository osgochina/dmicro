package proxy_test

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/proxy"
	"github.com/osgochina/dmicro/logger"
	"strconv"
	"testing"
	"time"
)

type Request struct {
	One int
	Two int
}

type Response struct {
	Three int
}

type math struct{ drpc.CallCtx }

func (m *math) Add(arg *Request) (*Response, *drpc.Status) {
	return &Response{Three: arg.One + arg.Two}, nil
}

type mathPush struct{ drpc.PushCtx }

func (m *mathPush) Push(arg *Request) *drpc.Status {
	return nil
}

func newSession(t *gtest.T, newProxy func() drpc.Plugin) drpc.Session {
	srv := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	})
	srv.RouteCall(new(math))
	srv.RoutePush(new(mathPush))
	go srv.ListenAndServe()
	time.Sleep(time.Second)
	srv1 := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  8080,
		PrintDetail: true,
	},
		newProxy(),
	)
	go srv1.ListenAndServe()
	time.Sleep(time.Second)
	cli := drpc.NewEndpoint(drpc.EndpointConfig{
		PrintDetail: true,
	})
	sess, stat := cli.Dial(":" + strconv.Itoa(8080))
	if !stat.OK() {
		t.Fatal(stat)
	}
	return sess
}

func newUnknownProxy() drpc.Plugin {
	cli := drpc.NewEndpoint(drpc.EndpointConfig{RedialTimes: 3})
	var sess drpc.Session
	var stat *drpc.Status
DIAL:
	sess, stat = cli.Dial(":9090")
	if !stat.OK() {
		logger.Warningf("%v", stat)
		time.Sleep(time.Second * 3)
		goto DIAL
	}
	return proxy.NewProxyPlugin(func(label *proxy.Label) proxy.Forwarder {
		logger.Infof("label RealIP:%s", label.RealIP)
		logger.Infof("label SessionID:%s", label.SessionID)
		logger.Infof("label ServiceMethod:%s", label.ServiceMethod)
		return sess
	})
}

func newUnknownCallProxy() drpc.Plugin {
	cli := drpc.NewEndpoint(drpc.EndpointConfig{RedialTimes: 3})
	var sess drpc.Session
	var stat *drpc.Status
DIAL:
	sess, stat = cli.Dial(":9090")
	if !stat.OK() {
		logger.Warningf("%v", stat)
		time.Sleep(time.Second * 3)
		goto DIAL
	}
	return proxy.NewProxyCallPlugin(func(label *proxy.Label) proxy.CallForwarder {
		logger.Infof("label RealIP:%s", label.RealIP)
		logger.Infof("label SessionID:%s", label.SessionID)
		logger.Infof("label ServiceMethod:%s", label.ServiceMethod)
		return sess
	})
}

func newUnknownPushProxy() drpc.Plugin {
	cli := drpc.NewEndpoint(drpc.EndpointConfig{RedialTimes: 3})
	var sess drpc.Session
	var stat *drpc.Status
DIAL:
	sess, stat = cli.Dial(":9090")
	if !stat.OK() {
		logger.Warningf("%v", stat)
		time.Sleep(time.Second * 3)
		goto DIAL
	}
	return proxy.NewProxyPushPlugin(func(label *proxy.Label) proxy.PushForwarder {
		logger.Infof("label RealIP:%s", label.RealIP)
		logger.Infof("label SessionID:%s", label.SessionID)
		logger.Infof("label ServiceMethod:%s", label.ServiceMethod)
		return sess
	})
}

func TestProxy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sess := newSession(t, newUnknownProxy)
		var result Response
		stat := sess.Call(
			"/math/add",
			&Request{One: 1, Two: 2},
			&result,
		).Status()
		t.Assert(stat.OK(), true)
		t.Assert(result.Three, 3)
		t.Logf("测试proxy：1+2=%d", result.Three)
		stat2 := sess.Push(
			"/math_push/push",
			&Request{One: 1, Two: 2},
		)
		t.Assert(stat2.OK(), true)
		time.Sleep(1 * time.Second)
	})
}

func TestCallProxy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sess := newSession(t, newUnknownCallProxy)
		var result Response
		stat := sess.Call(
			"/math/add",
			&Request{One: 1, Two: 2},
			&result,
		).Status()
		t.Assert(stat.OK(), true)
		t.Assert(result.Three, 3)
		t.Logf("测试proxy：1+2=%d", result.Three)
	})
}

func TestPushProxy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sess := newSession(t, newUnknownPushProxy)
		stat2 := sess.Push(
			"/math_push/push",
			&Request{One: 1, Two: 2},
		)
		t.Assert(stat2.OK(), true)
		time.Sleep(1 * time.Second)
	})
}
