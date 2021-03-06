package drpc_test

import (
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc"
	"net"
	"testing"
	"time"
)

type EventTestPlugin struct {
	BeforeNewEndpointFunc    func(cfg *drpc.EndpointConfig, plugins *drpc.PluginContainer) error
	AfterNewEndpointFunc     func(e drpc.EarlyEndpoint) error
	BeforeCloseEndpointFunc  func(drpc.Endpoint) error
	AfterCloseEndpointFunc   func(drpc.Endpoint, error) error
	AfterRegRouterFunc       func(*drpc.Handler) error
	AfterListenFunc          func(net.Addr) error
	BeforeDialFunc           func(addr string, isRedial bool) *drpc.Status
	AfterDialFunc            func(sess drpc.EarlySession, isRedial bool) *drpc.Status
	AfterDialFailFunc        func(sess drpc.EarlySession, err error, isRedial bool) *drpc.Status
	AfterAcceptFunc          func(sess drpc.EarlySession) *drpc.Status
	BeforeWriteCallFunc      func(ctx drpc.WriteCtx) *drpc.Status
	AfterWriteCallFunc       func(ctx drpc.WriteCtx) *drpc.Status
	BeforeWriteReplyFunc     func(ctx drpc.WriteCtx) *drpc.Status
	AfterWriteReplyFunc      func(ctx drpc.WriteCtx) *drpc.Status
	BeforeWritePushFunc      func(ctx drpc.WriteCtx) *drpc.Status
	AfterWritePushFunc       func(ctx drpc.WriteCtx) *drpc.Status
	BeforeReadHeaderFunc     func(ctx drpc.EarlyCtx) error
	AfterReadCallHeaderFunc  func(ctx drpc.ReadCtx) *drpc.Status
	BeforeReadCallBodyFunc   func(ctx drpc.ReadCtx) *drpc.Status
	AfterReadCallBodyFunc    func(ctx drpc.ReadCtx) *drpc.Status
	AfterReadPushHeaderFunc  func(ctx drpc.ReadCtx) *drpc.Status
	BeforeReadPushBodyFunc   func(ctx drpc.ReadCtx) *drpc.Status
	AfterReadPushBodyFunc    func(ctx drpc.ReadCtx) *drpc.Status
	AfterReadReplyHeaderFunc func(ctx drpc.ReadCtx) *drpc.Status
	BeforeReadReplyBodyFunc  func(ctx drpc.ReadCtx) *drpc.Status
	AfterReadReplyBodyFunc   func(ctx drpc.ReadCtx) *drpc.Status
	AfterDisconnectFunc      func(sess drpc.BaseSession) *drpc.Status
}

func NewEventTestPlugin() *EventTestPlugin {
	return &EventTestPlugin{}
}
func (that *EventTestPlugin) Name() string {
	return "EventTestPlugin"
}
func (that *EventTestPlugin) BeforeNewEndpoint(cfg *drpc.EndpointConfig, plugins *drpc.PluginContainer) error {
	if that.BeforeNewEndpointFunc == nil {
		return nil
	}
	return that.BeforeNewEndpointFunc(cfg, plugins)
}

func (that *EventTestPlugin) AfterNewEndpoint(e drpc.EarlyEndpoint) error {
	if that.AfterNewEndpointFunc == nil {
		return nil
	}
	return that.AfterNewEndpointFunc(e)
}

func (that *EventTestPlugin) BeforeCloseEndpoint(e drpc.Endpoint) error {
	if that.BeforeCloseEndpointFunc == nil {
		return nil
	}
	return that.BeforeCloseEndpointFunc(e)
}

func (that *EventTestPlugin) AfterCloseEndpoint(e drpc.Endpoint, err error) error {
	if that.AfterCloseEndpointFunc == nil {
		return nil
	}
	return that.AfterCloseEndpointFunc(e, err)
}

func (that *EventTestPlugin) AfterRegRouter(h *drpc.Handler) error {
	if that.AfterRegRouterFunc == nil {
		return nil
	}
	return that.AfterRegRouterFunc(h)
}

func (that *EventTestPlugin) AfterListen(addr net.Addr) error {
	if that.AfterListenFunc == nil {
		return nil
	}
	return that.AfterListenFunc(addr)
}

func (that *EventTestPlugin) BeforeDial(addr string, isRedial bool) *drpc.Status {
	if that.BeforeDialFunc == nil {
		return nil
	}
	return that.BeforeDialFunc(addr, isRedial)
}

func (that *EventTestPlugin) AfterDial(sess drpc.EarlySession, isRedial bool) *drpc.Status {
	if that.AfterDialFunc == nil {
		return nil
	}
	return that.AfterDialFunc(sess, isRedial)
}

func (that *EventTestPlugin) AfterDialFail(sess drpc.EarlySession, err error, isRedial bool) *drpc.Status {
	if that.AfterDialFailFunc == nil {
		return nil
	}
	return that.AfterDialFailFunc(sess, err, isRedial)
}

func (that *EventTestPlugin) AfterAccept(sess drpc.EarlySession) *drpc.Status {
	if that.AfterAcceptFunc == nil {
		return nil
	}
	return that.AfterAcceptFunc(sess)
}

func (that *EventTestPlugin) BeforeWriteCall(ctx drpc.WriteCtx) *drpc.Status {
	if that.BeforeWriteCallFunc == nil {
		return nil
	}
	return that.BeforeWriteCallFunc(ctx)
}

func (that *EventTestPlugin) AfterWriteCall(ctx drpc.WriteCtx) *drpc.Status {
	if that.AfterWriteCallFunc == nil {
		return nil
	}
	return that.AfterWriteCallFunc(ctx)
}

func (that *EventTestPlugin) BeforeWriteReply(ctx drpc.WriteCtx) *drpc.Status {
	if that.BeforeWriteReplyFunc == nil {
		return nil
	}
	return that.BeforeWriteReplyFunc(ctx)
}

func (that *EventTestPlugin) AfterWriteReply(ctx drpc.WriteCtx) *drpc.Status {
	if that.AfterWriteReplyFunc == nil {
		return nil
	}
	return that.AfterWriteReplyFunc(ctx)
}

func (that *EventTestPlugin) BeforeWritePush(ctx drpc.WriteCtx) *drpc.Status {
	if that.BeforeWritePushFunc == nil {
		return nil
	}
	return that.BeforeWritePushFunc(ctx)
}

func (that *EventTestPlugin) AfterWritePush(ctx drpc.WriteCtx) *drpc.Status {
	if that.AfterWritePushFunc == nil {
		return nil
	}
	return that.AfterWritePushFunc(ctx)
}

func (that *EventTestPlugin) BeforeReadHeader(ctx drpc.EarlyCtx) error {
	if that.BeforeReadHeaderFunc == nil {
		return nil
	}
	return that.BeforeReadHeaderFunc(ctx)
}

func (that *EventTestPlugin) AfterReadCallHeader(ctx drpc.ReadCtx) *drpc.Status {
	if that.AfterReadCallHeaderFunc == nil {
		return nil
	}
	return that.AfterReadCallHeaderFunc(ctx)
}

func (that *EventTestPlugin) BeforeReadCallBody(ctx drpc.ReadCtx) *drpc.Status {
	if that.BeforeReadCallBodyFunc == nil {
		return nil
	}
	return that.BeforeReadCallBodyFunc(ctx)
}

func (that *EventTestPlugin) AfterReadCallBody(ctx drpc.ReadCtx) *drpc.Status {
	if that.AfterReadCallBodyFunc == nil {
		return nil
	}
	return that.AfterReadCallBodyFunc(ctx)
}

func (that *EventTestPlugin) AfterReadPushHeader(ctx drpc.ReadCtx) *drpc.Status {
	if that.AfterReadPushHeaderFunc == nil {
		return nil
	}
	return that.AfterReadPushHeaderFunc(ctx)
}

func (that *EventTestPlugin) BeforeReadPushBody(ctx drpc.ReadCtx) *drpc.Status {
	if that.BeforeReadPushBodyFunc == nil {
		return nil
	}
	return that.BeforeReadPushBodyFunc(ctx)
}

func (that *EventTestPlugin) AfterReadPushBody(ctx drpc.ReadCtx) *drpc.Status {
	if that.AfterReadPushBodyFunc == nil {
		return nil
	}
	return that.AfterReadPushBodyFunc(ctx)
}

func (that *EventTestPlugin) AfterReadReplyHeader(ctx drpc.ReadCtx) *drpc.Status {
	if that.AfterReadReplyHeaderFunc == nil {
		return nil
	}
	return that.AfterReadReplyHeaderFunc(ctx)
}

func (that *EventTestPlugin) BeforeReadReplyBody(ctx drpc.ReadCtx) *drpc.Status {
	if that.BeforeReadReplyBodyFunc == nil {
		return nil
	}
	return that.BeforeReadReplyBodyFunc(ctx)
}

func (that *EventTestPlugin) AfterReadReplyBody(ctx drpc.ReadCtx) *drpc.Status {
	if that.AfterReadReplyBodyFunc == nil {
		return nil
	}
	return that.AfterReadReplyBodyFunc(ctx)
}

func (that *EventTestPlugin) AfterDisconnect(sess drpc.BaseSession) *drpc.Status {
	if that.AfterDisconnectFunc == nil {
		return nil
	}
	return that.AfterDisconnectFunc(sess)
}

func TestEventPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plugins := NewEventTestPlugin()
		plugins.BeforeNewEndpointFunc = func(cfg *drpc.EndpointConfig, plugins *drpc.PluginContainer) error {
			t.Log("???????????????BeforeNewEndpoint")
			//t.Assert(cfg.ListenIP, "127.0.0.1")
			//t.Assert(cfg.ListenPort, 9901)
			//t.Assert(cfg.Network, "tcp")
			//t.Assert(plugins.GetByName("EventTestPlugin").Name(), "EventTestPlugin")
			//t.AssertEQ(len(plugins.GetAll()), 1)
			return nil
		}

		plugins.AfterNewEndpointFunc = func(e drpc.EarlyEndpoint) error {
			t.Log("???????????????AfterNewEndpoint")
			//t.AssertEQ(e.CountSession(), 0)
			//t.AssertEQ(len(e.PluginContainer().GetAll()), 1)
			return nil
		}

		plugins.BeforeCloseEndpointFunc = func(e drpc.Endpoint) error {
			t.Log("???????????????BeforeCloseEndpoint")
			//t.AssertNE(e.Router(), nil)
			return nil
		}

		plugins.AfterCloseEndpointFunc = func(e drpc.Endpoint, err error) error {
			t.Log("???????????????AfterCloseEndpoint")
			//t.AssertNE(e.Router(), nil)
			//t.Assert(err, nil)
			return nil
		}

		plugins.AfterRegRouterFunc = func(h *drpc.Handler) error {
			t.Log("???????????????AfterRegRouter")
			//t.AssertEQ(h.RouterTypeName(), "CALL")
			//t.AssertEQ(h.Name(), "/math/add")
			return nil
		}

		plugins.AfterListenFunc = func(addr net.Addr) error {
			t.Log("???????????????AfterListen")
			//t.AssertEQ(addr.String(), "127.0.0.1:9901")
			return nil
		}

		plugins.BeforeDialFunc = func(addr string, isRedial bool) *drpc.Status {
			t.Log("???????????????BeforeDial")
			return nil
		}

		plugins.AfterDialFunc = func(sess drpc.EarlySession, isRedial bool) *drpc.Status {
			t.Log("???????????????AfterDial")
			return nil
		}

		plugins.AfterDialFailFunc = func(sess drpc.EarlySession, err error, isRedial bool) *drpc.Status {
			t.Log("???????????????AfterDialFail")
			t.Log(err)
			return nil
		}

		plugins.AfterAcceptFunc = func(sess drpc.EarlySession) *drpc.Status {
			t.Log("???????????????AfterAccept")
			return nil
		}

		plugins.BeforeWriteCallFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Log("???????????????BeforeWriteCall")
			return nil
		}

		plugins.AfterWriteCallFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Log("???????????????AfterWriteCall")
			return nil
		}

		plugins.BeforeWriteReplyFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Log("???????????????BeforeWriteReply")
			return nil
		}

		plugins.AfterWriteReplyFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Log("???????????????AfterWriteReply")
			return nil
		}

		plugins.BeforeWritePushFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Log("???????????????BeforeWritePush")
			return nil
		}

		plugins.AfterWritePushFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Log("???????????????AfterWritePush")
			return nil
		}

		plugins.BeforeReadHeaderFunc = func(ctx drpc.EarlyCtx) error {
			t.Log("???????????????BeforeReadHeader")
			return nil
		}

		plugins.AfterReadCallHeaderFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????AfterReadCallHeader")
			return nil
		}
		plugins.BeforeReadCallBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????BeforeReadCallBody")
			return nil
		}

		plugins.AfterReadCallBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????AfterReadCallBody")
			return nil
		}

		plugins.AfterReadPushHeaderFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????AfterReadPushHeader")
			return nil
		}

		plugins.BeforeReadPushBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????BeforeReadPushBody")
			return nil
		}

		plugins.AfterReadPushBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????AfterReadPushBody")
			return nil
		}
		plugins.AfterReadReplyHeaderFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????AfterReadReplyHeader")
			return nil
		}

		plugins.BeforeReadReplyBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????BeforeReadReplyBody")
			return nil
		}

		plugins.AfterReadReplyBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Log("???????????????AfterReadReplyBody")
			return nil
		}

		plugins.AfterDisconnectFunc = func(sess drpc.BaseSession) *drpc.Status {
			t.Log("???????????????AfterDisconnect")
			return nil
		}

		endpointSvr := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:    "tcp",
			ListenIP:   "127.0.0.1",
			ListenPort: 9901,
		}, plugins)
		endpointSvr.RouteCall(new(Math))

		endpointCli := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:   "tcp",
			LocalIP:   "127.0.0.1",
			LocalPort: 9902,
		}, plugins)

		gtimer.AddOnce(10*time.Second, func() {
			_ = endpointCli.Close()
			_ = endpointSvr.Close()
		})
		gtimer.AddOnce(1*time.Second, func() {
			sess, status := endpointCli.Dial("127.0.0.1:9901")
			if !status.OK() {
				t.Fatal("dial 127.0.0.1:9901 fail")
			}
			t.Assert(sess.ID(), "127.0.0.1:9902")
			var result int
			status = sess.Call("/math/add", []int{1, 2, 3}, &result).Status()
			if !status.OK() {
				t.Fatalf("/math/add fail,%v", status)
			}
			t.Log(result)
			t.Assert(result, 6)
		})
		_ = endpointSvr.ListenAndServe()
		time.Sleep(1 * time.Second)
	})
}

func TestEndpointPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plugins := NewEventTestPlugin()
		plugins.BeforeNewEndpointFunc = func(cfg *drpc.EndpointConfig, plugins *drpc.PluginContainer) error {
			t.Log("???????????????BeforeNewEndpoint")
			t.Assert(cfg.ListenIP, "127.0.0.1")
			t.Assert(cfg.ListenPort, 9901)
			t.Assert(cfg.Network, "tcp")
			t.Assert(plugins.GetByName("EventTestPlugin").Name(), "EventTestPlugin")
			t.AssertEQ(len(plugins.GetAll()), 1)
			return nil
		}

		plugins.AfterNewEndpointFunc = func(e drpc.EarlyEndpoint) error {
			t.Log("???????????????AfterNewEndpoint")
			t.AssertEQ(e.CountSession(), 0)
			t.AssertEQ(len(e.PluginContainer().GetAll()), 1)
			return nil
		}

		plugins.BeforeCloseEndpointFunc = func(e drpc.Endpoint) error {
			t.Log("???????????????BeforeCloseEndpoint")
			t.AssertNE(e.Router(), nil)
			return nil
		}

		plugins.AfterCloseEndpointFunc = func(e drpc.Endpoint, err error) error {
			t.Log("???????????????AfterCloseEndpoint")
			t.AssertNE(e.Router(), nil)
			t.Assert(err, nil)
			return nil
		}
		endpoint := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:    "tcp",
			ListenIP:   "127.0.0.1",
			ListenPort: 9901,
		}, plugins)
		_ = endpoint.Close()
	})
}

func TestRegRouterPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plugins := NewEventTestPlugin()
		plugins.AfterRegRouterFunc = func(h *drpc.Handler) error {
			t.Logf("???????????????AfterRegRouter,type:%s,name:%s", h.RouterTypeName(), h.Name())
			return nil
		}
		endpoint := drpc.NewEndpoint(drpc.EndpointConfig{}, plugins)
		endpoint.RouteCall(new(Math))
		endpoint.RouteCallFunc((*Math).AddFunc)
		endpoint.RoutePush(new(MathPush))
		endpoint.RoutePushFunc((*MathPush).AddFunc)
	})
}

func TestListenPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plugins := NewEventTestPlugin()

		plugins.AfterListenFunc = func(addr net.Addr) error {
			t.Logf("???????????????AfterListen,%s", addr.String())
			return nil
		}

		plugins.BeforeDialFunc = func(addr string, isRedial bool) *drpc.Status {
			t.Logf("???????????????BeforeDial,%s", addr)
			return nil
		}

		plugins.AfterDialFunc = func(sess drpc.EarlySession, isRedial bool) *drpc.Status {
			t.Logf("???????????????AfterDial,remoteAddr:%s,localAddr:%s", sess.RemoteAddr().String(), sess.LocalAddr())
			return nil
		}

		plugins.AfterDialFailFunc = func(sess drpc.EarlySession, err error, isRedial bool) *drpc.Status {
			t.Logf("???????????????AfterDialFail,err:%v", err)
			return nil
		}

		plugins.AfterAcceptFunc = func(sess drpc.EarlySession) *drpc.Status {
			t.Logf("???????????????AfterAccept,remoteAddr:%s,localAddr:%s", sess.RemoteAddr().String(), sess.LocalAddr())
			return nil
		}
		plugins.AfterDisconnectFunc = func(sess drpc.BaseSession) *drpc.Status {
			t.Logf("???????????????AfterDisconnect,remoteAddr:%s,localAddr:%s", sess.RemoteAddr().String(), sess.LocalAddr())
			return nil
		}
		endpointSvr := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:    "tcp",
			ListenIP:   "127.0.0.1",
			ListenPort: 9901,
		}, plugins)

		endpointCli := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:   "tcp",
			LocalIP:   "127.0.0.1",
			LocalPort: 0,
		}, plugins)
		gtimer.AddOnce(3*time.Second, func() {
			_ = endpointCli.Close()
			_ = endpointSvr.Close()
		})
		gtimer.AddOnce(1*time.Second, func() {
			sess, status := endpointCli.Dial("127.0.0.1:9901")
			if !status.OK() {
				t.Fatal("dial 127.0.0.1:9901 fail")
			}
			t.Log(sess.ID())
		})
		_ = endpointSvr.ListenAndServe()
		time.Sleep(time.Second * 1)
	})
}

func TestCallPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plugins := NewEventTestPlugin()

		plugins.BeforeWriteCallFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Logf("???????????????BeforeWriteCall,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.AfterWriteCallFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Logf("???????????????AfterWriteCall,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())

			return nil
		}

		plugins.BeforeWriteReplyFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Logf("???????????????BeforeWriteReply,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.AfterWriteReplyFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Logf("???????????????AfterWriteReply,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.BeforeReadHeaderFunc = func(ctx drpc.EarlyCtx) error {
			t.Logf("???????????????BeforeReadHeader,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.AfterReadCallHeaderFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????AfterReadCallHeader,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}
		plugins.BeforeReadCallBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????BeforeReadCallBody,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.AfterReadCallBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????AfterReadCallBody,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.AfterReadReplyHeaderFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????AfterReadReplyHeader,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.BeforeReadReplyBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????BeforeReadReplyBody,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.AfterReadReplyBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????AfterReadReplyBody,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		endpointSvr := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:    "tcp",
			ListenIP:   "127.0.0.1",
			ListenPort: 9901,
		}, plugins)
		endpointSvr.RouteCall(new(Math))
		endpointCli := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:   "tcp",
			LocalIP:   "127.0.0.1",
			LocalPort: 0,
		}, plugins)
		gtimer.AddOnce(3*time.Second, func() {
			_ = endpointCli.Close()
			_ = endpointSvr.Close()
		})
		gtimer.AddOnce(1*time.Second, func() {
			sess, status := endpointCli.Dial("127.0.0.1:9901")
			if !status.OK() {
				t.Fatal("dial 127.0.0.1:9901 fail")
			}
			t.Log(sess.ID())
			var result int
			status = sess.Call("/math/add", []int{1, 2, 3}, &result).Status()
			if !status.OK() {
				t.Fatalf("/math/add fail,%v", status)
			}
			t.Log(result)
			t.Assert(result, 6)
		})
		_ = endpointSvr.ListenAndServe()
		time.Sleep(time.Second * 1)
	})
}

func TestPushPlugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plugins := NewEventTestPlugin()

		plugins.BeforeWritePushFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Logf("???????????????BeforeWritePush,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.AfterWritePushFunc = func(ctx drpc.WriteCtx) *drpc.Status {
			t.Logf("???????????????AfterWritePush,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())

			return nil
		}

		plugins.BeforeReadHeaderFunc = func(ctx drpc.EarlyCtx) error {
			t.Logf("???????????????BeforeReadHeader,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())
			return nil
		}

		plugins.AfterReadPushHeaderFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????AfterReadPushHeader,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())

			return nil
		}

		plugins.BeforeReadPushBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????BeforeReadPushBody,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())

			return nil
		}

		plugins.AfterReadPushBodyFunc = func(ctx drpc.ReadCtx) *drpc.Status {
			t.Logf("???????????????AfterReadPushBody,LocalAddr:%s,RemoteAddr:%s", ctx.Session().LocalAddr(), ctx.Session().RemoteAddr())

			return nil
		}

		endpointSvr := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:    "tcp",
			ListenIP:   "127.0.0.1",
			ListenPort: 9901,
		}, plugins)
		endpointSvr.RoutePush(new(MathPush))
		endpointCli := drpc.NewEndpoint(drpc.EndpointConfig{
			Network:   "tcp",
			LocalIP:   "127.0.0.1",
			LocalPort: 0,
		}, plugins)
		gtimer.AddOnce(3*time.Second, func() {
			_ = endpointCli.Close()
			_ = endpointSvr.Close()
		})
		gtimer.AddOnce(1*time.Second, func() {
			sess, status := endpointCli.Dial("127.0.0.1:9901")
			if !status.OK() {
				t.Fatal("dial 127.0.0.1:9901 fail")
			}
			t.Log(sess.ID())

			status = sess.Push("/math_push/add", []int{1, 2, 3})
			if !status.OK() {
				t.Fatalf("/math_push/add fail,%v", status)
			}
		})
		_ = endpointSvr.ListenAndServe()
		time.Sleep(time.Second * 1)
	})
}

type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}

func (m *Math) AddFunc(arg *[]int) (int, *drpc.Status) {
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}

type MathPush struct {
	drpc.PushCtx
}

func (m *MathPush) Add(arg *[]int) *drpc.Status {
	var r int
	for _, a := range *arg {
		r += a
	}
	return nil
}

func (m *MathPush) AddFunc(arg *[]int) *drpc.Status {
	var r int
	for _, a := range *arg {
		r += a
	}
	return nil
}
