package proxy_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/internal"
	"github.com/osgochina/dmicro/drpc/plugin/proxy"
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

func newSession(t *gtest.T, newProxy func() drpc.Plugin, backendPort int, proxyPort int) drpc.Session {
	srv := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  uint16(backendPort),
		PrintDetail: true,
	})
	srv.RouteCall(new(math))
	srv.RoutePush(new(mathPush))
	go srv.ListenAndServe()
	time.Sleep(time.Second)
	srv1 := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  uint16(proxyPort),
		PrintDetail: true,
	},
		newProxy(),
	)
	go srv1.ListenAndServe()
	time.Sleep(time.Second)
	cli := drpc.NewEndpoint(drpc.EndpointConfig{
		PrintDetail: true,
	})
	sess, stat := cli.Dial(":" + strconv.Itoa(proxyPort))
	if !stat.OK() {
		t.Fatal(stat)
	}
	return sess
}

func newUnknownProxy(backendPort int) drpc.Plugin {
	cli := drpc.NewEndpoint(drpc.EndpointConfig{RedialTimes: 3})
	var sess drpc.Session
	var stat *drpc.Status
DIAL:
	sess, stat = cli.Dial(":" + strconv.Itoa(backendPort))
	if !stat.OK() {
		internal.Warningf(context.TODO(), "%v", stat)
		time.Sleep(time.Second * 3)
		goto DIAL
	}
	return proxy.NewProxyPlugin(func(label *proxy.Label) proxy.Forwarder {
		internal.Infof(context.TODO(), "label RealIP:%s", label.RealIP)
		internal.Infof(context.TODO(), "label SessionID:%s", label.SessionID)
		internal.Infof(context.TODO(), "label ServiceMethod:%s", label.ServiceMethod)
		return sess
	})
}

func newUnknownCallProxy(backendPort int) drpc.Plugin {
	cli := drpc.NewEndpoint(drpc.EndpointConfig{RedialTimes: 3})
	var sess drpc.Session
	var stat *drpc.Status
DIAL:
	sess, stat = cli.Dial(":" + strconv.Itoa(backendPort))
	if !stat.OK() {
		internal.Warningf(context.TODO(), "%v", stat)
		time.Sleep(time.Second * 3)
		goto DIAL
	}
	return proxy.NewProxyCallPlugin(func(label *proxy.Label) proxy.CallForwarder {
		internal.Infof(context.TODO(), "label RealIP:%s", label.RealIP)
		internal.Infof(context.TODO(), "label SessionID:%s", label.SessionID)
		internal.Infof(context.TODO(), "label ServiceMethod:%s", label.ServiceMethod)
		return sess
	})
}

func newUnknownPushProxy(backendPort int) drpc.Plugin {
	cli := drpc.NewEndpoint(drpc.EndpointConfig{RedialTimes: 3})
	var sess drpc.Session
	var stat *drpc.Status
DIAL:
	sess, stat = cli.Dial(":" + strconv.Itoa(backendPort))
	if !stat.OK() {
		internal.Warningf(context.TODO(), "%v", stat)
		time.Sleep(time.Second * 3)
		goto DIAL
	}
	return proxy.NewProxyPushPlugin(func(label *proxy.Label) proxy.PushForwarder {
		internal.Infof(context.TODO(), "label RealIP:%s", label.RealIP)
		internal.Infof(context.TODO(), "label SessionID:%s", label.SessionID)
		internal.Infof(context.TODO(), "label ServiceMethod:%s", label.ServiceMethod)
		return sess
	})
}

func TestProxy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sess := newSession(t, func() drpc.Plugin { return newUnknownProxy(9099) }, 9099, 8080)
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
		sess := newSession(t, func() drpc.Plugin { return newUnknownCallProxy(9091) }, 9091, 8081)
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
		sess := newSession(t, func() drpc.Plugin { return newUnknownPushProxy(9092) }, 9092, 8082)
		stat2 := sess.Push(
			"/math_push/push",
			&Request{One: 1, Two: 2},
		)
		t.Assert(stat2.OK(), true)
		time.Sleep(1 * time.Second)
	})
}
