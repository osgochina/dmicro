package prometheus

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"testing"
)

func TestNewCounterVec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cVec := NewCounterVec(&CounterVecOpts{
			Namespace: "rpc_server",
			Subsystem: "call",
			Name:      "total",
			Help:      "rpc client call code total",
		})
		defer cVec.Close()
		t.AssertNE(cVec, nil)
		cVec = NewCounterVec(nil)
		t.Assert(cVec, nil)
	})
}

func TestCounterVec_Inc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cVec := NewCounterVec(&CounterVecOpts{
			Namespace: "rpc_server",
			Subsystem: "call",
			Name:      "total",
			Help:      "rpc client call code total",
			Labels:    []string{"path", "code"},
		})
		defer cVec.Close()
		cVec.Inc("/admin", "404")
		cVec.Inc("/admin", "404")
		c := testutil.ToFloat64(cVec.(*counterVec).counter)
		t.Assert(c, float64(2))
	})
}

func TestCounterVec_Add(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cVec := NewCounterVec(&CounterVecOpts{
			Namespace: "rpc_server",
			Subsystem: "call",
			Name:      "total",
			Help:      "rpc client call code total",
			Labels:    []string{"path", "code"},
		})
		defer cVec.Close()
		cVec.Inc("/admin", "404")
		cVec.Add(10, "/admin", "404")
		c := testutil.ToFloat64(cVec.(*counterVec).counter)
		t.Assert(c, float64(11))
	})
}
