package message

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc/status"
	"github.com/osgochina/dmicro/drpc/tfilter"
	"testing"
)

func TestMessageSizeLimit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Logf("MessageSizeLimit: %d", MsgSizeLimit())
		maxLen := 1024 * 1024 * 8
		SetMsgSizeLimit(uint32(maxLen))
		t.Assert(MsgSizeLimit(), maxLen)
	})
}

func TestMessageString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tfilter.RegGzip(5)
		var m = GetMessage()
		defer PutMessage(m)
		m.SetSeq(21)
		err := m.PipeTFilter().Append(tfilter.GzipId)
		t.Assert(err, nil)
		m.SetMType(3)
		err = m.SetSize(300)
		t.Assert(err, nil)
		m.SetBody(map[string]int{"a": 1})
		m.SetServiceMethod("service/method")
		m.SetBodyCodec(5)
		m.SetStatus(status.New(400, "this is msg", "this is cause"))
		m.Meta().Set("key", "value")
		t.Logf("%%s:%s", m.String())
		t.Logf("%%v:%v", m)
		t.Logf("%%#v:%#v", m)
		t.Logf("%%+v:%+v", m)
	})

}
