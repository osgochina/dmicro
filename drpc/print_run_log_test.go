package drpc

import (
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func TestEnablePrintRunLog(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_ = glog.SetLevelStr("ALL")
		t.Assert(enablePrintRunLog(), true)
		_ = glog.SetLevelStr("DEV")
		t.Assert(enablePrintRunLog(), true)
		_ = glog.SetLevelStr("DEVELOP")
		t.Assert(enablePrintRunLog(), true)
		_ = glog.SetLevelStr("PROD")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("PRODUCT")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("DEBU")
		t.Assert(enablePrintRunLog(), true)
		_ = glog.SetLevelStr("DEBUG")
		t.Assert(enablePrintRunLog(), true)
		_ = glog.SetLevelStr("INFO")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("NOTI")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("NOTICE")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("WARN")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("WARNING")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("ERRO")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("ERROR")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("CRIT")
		t.Assert(enablePrintRunLog(), false)
		_ = glog.SetLevelStr("CRITICAL")
		t.Assert(enablePrintRunLog(), false)
	})
}
