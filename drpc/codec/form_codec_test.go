package codec

import (
	"github.com/gogf/gf/test/gtest"
	"net/url"
	"testing"
)

func TestForm(t *testing.T) {
	var (
		data = "a=aaa&b=2&c=ccc1&c=ccc2&d=4&d=44"
		c    = new(FormCodec)
	)
	gtest.C(t, func(t *gtest.T) {
		var v1 url.Values
		err := c.Unmarshal([]byte(data), &v1)
		t.Assert(err, nil)
		t.Logf("Unmarshal 1: %#v", v1)

		var v2 map[string][]string
		err = c.Unmarshal([]byte(data), &v2)
		t.Assert(err, nil)
		t.Logf("Unmarshal 2: %#v", v2)

		type T struct {
			A string   `form:"a"`
			B int32    `form:"b"`
			C []string `form:"c"`
			D [2]int   `form:"d"`
		}
		var v3 T
		err = c.Unmarshal([]byte(data), &v3)
		t.Assert(err, nil)
		t.Logf("Unmarshal 3: %#v", v3)

		b3, err := c.Marshal(v3)
		t.Assert(err, nil)
		t.Logf("Marshal 3: %s", b3)
		var v4 interface{}
		err = c.Unmarshal(b3, &v4)
		t.Assert(err, nil)
		t.Logf("Unmarshal 4: %#v", v4)
	})

}
