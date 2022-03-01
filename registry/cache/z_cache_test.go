package cache

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func TestCache_GetService(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := New(nil)
		s, err := c.GetService("test")
		t.Assert(err, nil)
		t.Log(s)
	})
}
