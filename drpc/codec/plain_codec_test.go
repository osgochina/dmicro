package codec

import (
	"github.com/gogf/gf/v2/test/gtest"
	"reflect"
	"testing"
)

func TestPlain(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type (
			S string
			I int32
		)
		var (
			bytes = make([]byte, 0, 100)
			a     = []interface{}{S("sss"), I(123), []byte("1234567890"), []byte("asdfghjkl")}
			b     = []interface{}{new(S), new(I), make([]byte, 10), &bytes}
			c     = new(PlainCodec)
		)
		for k, v := range a {
			data, err := c.Marshal(v)
			t.Assert(err, nil)
			err = c.Unmarshal(data, b[k])
			t.Assert(err, nil)
			if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(b[k])).Interface(), v) {
				t.Logf("get: %v, but expect: %v", b[k], v)
			}
		}
	})
}
