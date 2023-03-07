package drpc

import (
	"fmt"
	"github.com/gogf/gf/v2/test/gtest"
	"reflect"
	"testing"
)

func TestHTTPServiceMethodMapper(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(HTTPServiceMethodMapper("AaBb", "AaBb"), "/AaBb/aa_bb")
		t.Assert(HTTPServiceMethodMapper("AaBb", "ABcXYz"), "/AaBb/abc_xyz")
		t.Assert(HTTPServiceMethodMapper("AaBb", "Aa__Bb"), "/AaBb/aa_bb")
		t.Assert(HTTPServiceMethodMapper("AaBb", "aa__bb"), "/AaBb/aa_bb")
		t.Assert(HTTPServiceMethodMapper("AaBb", "ABC__XYZ"), "/AaBb/abc_xyz")
		t.Assert(HTTPServiceMethodMapper("AaBb", "Aa_Bb"), "/AaBb/aa/bb")
		t.Assert(HTTPServiceMethodMapper("AaBb", "aa_bb"), "/AaBb/aa/bb")
		t.Assert(HTTPServiceMethodMapper("AaBb", "ABC_XYZ"), "/AaBb/abc/xyz")
		t.Assert(HTTPServiceMethodMapper("Abc", "Efg"), "/Abc/efg")
		t.Assert(HTTPServiceMethodMapper("Abc", "Efg_abc"), "/Abc/efg/abc")
		t.Assert(HTTPServiceMethodMapper("", ""), "/")
	})
}

func TestRPCServiceMethodMapper(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(RPCServiceMethodMapper("AaBb", "AaBb"), "AaBb.AaBb")
		t.Assert(RPCServiceMethodMapper("AaBb", "ABcXYz"), "AaBb.ABcXYz")
		t.Assert(RPCServiceMethodMapper("AaBb", "Aa__Bb"), "AaBb.Aa_Bb")
		t.Assert(RPCServiceMethodMapper("AaBb", "aa__bb"), "AaBb.aa_bb")
		t.Assert(RPCServiceMethodMapper("AaBb", "ABC__XYZ"), "AaBb.ABC_XYZ")
		t.Assert(RPCServiceMethodMapper("AaBb", "Aa_Bb"), "AaBb.Aa.Bb")
		t.Assert(RPCServiceMethodMapper("AaBb", "aa_bb"), "AaBb.aa.bb")
		t.Assert(RPCServiceMethodMapper("AaBb", "ABC_XYZ"), "AaBb.ABC.XYZ")
		t.Assert(RPCServiceMethodMapper("", ""), "")
	})
}

func TestNewRouter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		pluginContainer := newPluginContainer()
		router := newRouter(pluginContainer)
		names := router.RouteCall(new(Math))
		t.AssertIN("/math/add", names)
		name := router.RouteCallFunc((*Math).AddFunc)
		t.Assert(name, "/add_func")
		names = router.RoutePush(new(MathPush))
		t.AssertIN("/math_push/add", names)
		name = router.RoutePushFunc((*MathPush).AddFunc)
		t.Assert(name, "/add_func")
	})
}

func TestGroupRouter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		pluginContainer := newPluginContainer()
		router := newRouter(pluginContainer)
		group := router.SubRoute("github")
		names := group.RouteCall(new(Math))
		t.AssertIN("/github/math/add", names)
		name := group.RouteCallFunc((*Math).AddFunc)
		t.Assert(name, "/github/add_func")
		names = group.RoutePush(new(MathPush))
		t.AssertIN("/github/math_push/add", names)
		name = group.RoutePushFunc((*MathPush).AddFunc)
		t.Assert(name, "/github/add_func")
		router = group.Root()
		names = router.RouteCall(new(Math))
		t.AssertIN("/math/add", names)
	})
}

func TestRouteCall(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := newRouter(newPluginContainer())
		root.RouteCall(new(Math))
		h, found := root.subRouter.getCall("/math/add")
		t.Assert(found, true)
		h.handleFunc(newReadHandleCtx(), h.NewArgValue())
		t.Assert(h.RouterTypeName(), "CALL")
		t.Assert(h.Name(), "/math/add")
		t.Assert(h.IsCall(), true)
		var aa []int
		t.Assert(h.ArgElemType(), reflect.TypeOf(aa))
	})
}

type Math struct {
	CallCtx
}

func (m *Math) Add(arg *[]int) (int, *Status) {
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}

func (m *Math) AddFunc(arg *[]int) (int, *Status) {
	var r int
	for _, a := range *arg {
		r += a
	}
	return r, nil
}

type MathPush struct {
	PushCtx
}

func (m *MathPush) Add(arg *[]int) *Status {
	var r int
	for _, a := range *arg {
		r += a
	}
	fmt.Println(r)
	return nil
}

func (m *MathPush) AddFunc(arg *[]int) *Status {
	var r int
	for _, a := range *arg {
		r += a
	}
	return nil
}
