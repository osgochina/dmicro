package config

import (
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func TestEntry_Load(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		teatPath := gdebug.TestDataPath()
		config := NewConfig(teatPath + "/test.conf")
		err := config.Load()
		t.Assert(err, nil)
		entrys := config.GetPrograms()
		t.Assert(len(entrys), 6)
	})
}
