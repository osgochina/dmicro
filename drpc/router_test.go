package drpc

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc/status"
	"github.com/osgochina/dmicro/logger"
	"testing"
)

//func TestHTTPServiceMethodMapper(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		t.Assert(globalServiceMethodMapper("Abc", "Efg"), "/Abc/efg")
//		t.Assert(globalServiceMethodMapper("", ""), "/")
//	})
//}

func TestRouteCall(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := newRouter(newPluginContainer())
		root.RouteCall(new(Math))
		h, found := root.subRouter.getCall("/math/add")
		if found {
			h.handleFunc(newReadHandleCtx(), h.NewArgValue())
		}
	})
}

type Math struct {
	name string
	CallCtx
}

// Add handles addition request
func (m *Math) Add(arg *[]int) (int, *status.Status) {
	// test meta
	logger.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
