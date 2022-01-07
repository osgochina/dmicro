package websocket_test

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/mixer/websocket"
	"github.com/osgochina/dmicro/drpc/mixer/websocket/pbSubProto"
	"github.com/osgochina/dmicro/utils"
	"net/http"
	"testing"
	"time"
)

type Arg struct {
	A int
	B int `param:"<range:1:>"`
}

type P struct {
	drpc.CallCtx
}

func (p *P) Divide(arg *Arg) (int, *drpc.Status) {
	return arg.A / arg.B, nil
}

// 测试基本json协议
func TestJSONWebsocket(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := websocket.NewServer("/", drpc.EndpointConfig{ListenPort: 9090})
		srv.RouteCall(new(P))
		go srv.ListenAndServe()
		time.Sleep(time.Second * 1)
		cli := websocket.NewClient("/", drpc.EndpointConfig{})
		sess, stat := cli.Dial(":9090")
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

func TestJSONWebsocketTLS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := websocket.NewServer("/", drpc.EndpointConfig{ListenPort: 9090})
		srv.RouteCall(new(P))
		srv.SetTLSConfig(utils.GenerateTLSConfigForServer())
		go srv.ListenAndServeJSON()
		time.Sleep(time.Second * 1)
		cli := websocket.NewClient("/", drpc.EndpointConfig{})
		cli.SetTLSConfig(utils.GenerateTLSConfigForClient())
		sess, stat := cli.Dial(":9090")
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

func TestPbWebsocket(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := websocket.NewServer("/abc", drpc.EndpointConfig{ListenPort: 9091, PrintDetail: true})
		srv.RouteCall(new(P))
		go srv.ListenAndServeProtobuf()
		time.Sleep(time.Second * 1)
		cli := websocket.NewClient("/abc", drpc.EndpointConfig{PrintDetail: true})
		sess, err := cli.DialProtobuf(":9091")
		if err != nil {
			t.Fatal(err)
		}
		var result int
		stat := sess.Call("/p/divide", &Arg{
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

func TestPbWebsocketTLS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := websocket.NewServer("/abc", drpc.EndpointConfig{ListenPort: 9091, PrintDetail: true})
		srv.RouteCall(new(P))
		srv.SetTLSConfig(utils.GenerateTLSConfigForServer())
		go srv.ListenAndServeProtobuf()
		time.Sleep(time.Second * 1)
		cli := websocket.NewClient("/abc", drpc.EndpointConfig{PrintDetail: true})
		cli.SetTLSConfig(utils.GenerateTLSConfigForClient())
		sess, err := cli.DialProtobuf(":9091")
		if err != nil {
			t.Fatal(err)
		}
		var result int
		stat := sess.Call("/p/divide", &Arg{
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

func TestCustomizedWebsocket(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srv := drpc.NewEndpoint(drpc.EndpointConfig{})
		http.Handle("/ws", websocket.NewPbServeHandler(srv, nil))
		go http.ListenAndServe(":9092", nil)
		srv.RouteCall(new(P))
		time.Sleep(time.Second * 1)

		cli := drpc.NewEndpoint(drpc.EndpointConfig{}, websocket.NewDialPlugin("/ws"))
		sess, stat := cli.Dial(":9092", pbSubProto.NewPbSubProtoFunc())
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
