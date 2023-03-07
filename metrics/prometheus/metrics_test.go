package prometheus

import (
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/osgochina/dmicro/metrics"
	"testing"
)

type testPlugin struct {
}

func (that *testPlugin) Name() string {
	return "testPlugin"
}
func TestNewPromMetrics(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		pm := NewPromMetrics()
		t.Assert(pm.Options().Host, "0.0.0.0")
		t.Assert(pm.Options().Port, 9101)
		t.Assert(pm.Options().Path, "/metrics")
		pm = NewPromMetrics(
			metrics.OptPath("/metric"),
			metrics.OptServiceName("test"),
			metrics.OptPort(9999),
			metrics.OptHost("127.0.0.1"),
		)
		t.Assert(pm.Options().Host, "127.0.0.1")
		t.Assert(pm.Options().Path, "/metric")
		t.Assert(pm.Options().ServiceName, "test")
		t.Assert(pm.Options().Port, 9999)
		t.Assert(len(pm.Options().Plugins), 1)
	})
}

func TestNewPromMetrics_Plugin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		pm := NewPromMetrics(
			metrics.OptPlugin(new(testPlugin)),
		)
		t.Assert(len(pm.Options().Plugins), 2)
	})
}
