package prometheus

import (
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"strings"
	"testing"
)

func TestNewHistogramVec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		hVec := NewHistogramVec(&HistogramVecOpts{
			Name:    "duration_ms",
			Help:    "rpc server call duration(ms).",
			Buckets: []float64{1, 2, 3},
		})
		defer hVec.Close()

		t.AssertNE(hVec, nil)
		t.Assert(nil, NewHistogramVec(nil))
	})
}

func TestHistogramVec_Observe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		hVec := NewHistogramVec(&HistogramVecOpts{
			Name:    "counts",
			Help:    "rpc server call duration(ms).",
			Buckets: []float64{1, 2, 3},
			Labels:  []string{"method"},
		})
		defer hVec.Close()
		hVec.Observe(3, "/admin")
		metadata := `
		# HELP counts rpc server call duration(ms).
        # TYPE counts histogram
`
		val := `
		counts_bucket{method="/admin",le="1"} 0
		counts_bucket{method="/admin",le="2"} 0
		counts_bucket{method="/admin",le="3"} 1
		counts_bucket{method="/admin",le="+Inf"} 1
		counts_sum{method="/admin"} 3
        counts_count{method="/admin"} 1
`
		err := testutil.CollectAndCompare(hVec.(*histogramVec).histogram, strings.NewReader(metadata+val))
		t.Assert(err, nil)
	})
}
