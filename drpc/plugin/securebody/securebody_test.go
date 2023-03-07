package securebody_test

import (
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/plugin/securebody"
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

func newSession(t *gtest.T, port uint16) drpc.Session {
	p := securebody.NewSecureBodyPlugin("cipherkey1234567")
	srv := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  port,
		PrintDetail: true,
	})
	srv.RouteCall(new(math), p)
	go srv.ListenAndServe()
	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(drpc.EndpointConfig{
		PrintDetail: true,
	}, p)
	sess, stat := cli.Dial(":" + strconv.Itoa(int(port)))
	if !stat.OK() {
		t.Fatal(stat)
	}
	return sess
}

func TestSecureBodyPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sess := newSession(t, 9090)
		var result Response
		stat := sess.Call(
			"/math/add",
			&Request{One: 1, Two: 2},
			&result,
			securebody.WithSecureMeta(),
		).Status()
		t.Assert(stat.OK(), true)
		t.Assert(result.Three, 3)
		t.Logf("测试加密：1+2=%d", result.Three)
	})
}

func TestReplySecureBodyPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sess := newSession(t, 9090)
		var result Response
		stat := sess.Call(
			"/math/add",
			&Request{One: 1, Two: 2},
			&result,
			securebody.WithReplySecureMeta(true),
		).Status()
		t.Assert(stat.OK(), true)
		t.Assert(result.Three, 3)
		t.Logf("测试加密：1+2=%d", result.Three)
	})
}
