package dserver

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestGetGFConf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := "/tmp/config.json"
		err := gfile.PutContents(path, "{\n    \"addr\": \"127.0.0.1:80\"\n}")
		defer func() {
			_ = gfile.Remove(path)
		}()
		t.Assert(err, nil)
		cfg := newDServer("test").getGFConf(path)
		t.Assert(cfg.MustGet(context.TODO(), "addr").String(), "127.0.0.1:80")
	})
}
