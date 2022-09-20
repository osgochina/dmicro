package prometheus

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"testing"
)

func TestNewGaugeVec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gVec := NewGaugeVec(&GaugeVecOpts{
			Namespace: "rpc_server",
			Subsystem: "call",
			Name:      "duration",
			Help:      "rpc server call duration(ms).",
		})
		defer gVec.Close()
		t.AssertNE(gVec, nil)
		t.Assert(nil, NewGaugeVec(nil))
	})
}

func TestGaugeVec_Inc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gVec := NewGaugeVec(&GaugeVecOpts{
			Namespace: "rpc_server",
			Subsystem: "call",
			Name:      "duration_inc",
			Help:      "rpc server call duration(ms).",
			Labels:    []string{"path"},
		})
		defer gVec.Close()
		gVec.Inc("/admin")
		gVec.Inc("/admin")
		g := testutil.ToFloat64(gVec.(*gaugeVec).gauge)
		t.Assert(float64(2), g)
	})
}

func TestGaugeVec_Add(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gVec := NewGaugeVec(&GaugeVecOpts{
			Namespace: "rpc_server",
			Subsystem: "call",
			Name:      "duration_add",
			Help:      "rpc server call duration(ms).",
			Labels:    []string{"path"},
		})
		defer gVec.Close()
		gVec.Add(10, "/admin")
		gVec.Add(-5, "/admin")
		g := testutil.ToFloat64(gVec.(*gaugeVec).gauge)
		t.Assert(float64(5), g)
	})
}

func TestGaugeVec_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gVec := NewGaugeVec(&GaugeVecOpts{
			Namespace: "rpc_server",
			Subsystem: "call",
			Name:      "duration_set",
			Help:      "rpc server call duration(ms).",
			Labels:    []string{"path"},
		})
		defer gVec.Close()
		gVec.Set(99, "/admin")
		g := testutil.ToFloat64(gVec.(*gaugeVec).gauge)
		t.Assert(float64(99), g)
	})
}
