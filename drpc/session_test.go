package drpc

import (
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/drpc/proto"
	"testing"
)

func TestSessionNew(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		e := NewEndpoint(EndpointConfig{})
		var p []proto.ProtoFunc
		sess := newSession(e.(*endpoint), nil, p)
		t.Assert(sess.getStatus(), statusPreparing)
		sess.changeStatus(statusOk)
		t.Assert(sess.getStatus(), statusOk)

		sess.changeStatus(statusActiveClosing)
		t.Assert(sess.getStatus(), statusActiveClosing)
		t.Assert(sess.checkStatus(statusActiveClosing), true)
	})
}
