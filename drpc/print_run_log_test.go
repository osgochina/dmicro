package drpc

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/drpc/internal"
	"testing"
)

func TestEnablePrintRunLog(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_ = internal.SetLevelStr("ALL")
		t.Assert(enablePrintRunLog(), true)
		_ = internal.SetLevelStr("DEV")
		t.Assert(enablePrintRunLog(), true)
		_ = internal.SetLevelStr("DEVELOP")
		t.Assert(enablePrintRunLog(), true)
		_ = internal.SetLevelStr("PROD")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("PRODUCT")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("DEBU")
		t.Assert(enablePrintRunLog(), true)
		_ = internal.SetLevelStr("DEBUG")
		t.Assert(enablePrintRunLog(), true)
		_ = internal.SetLevelStr("INFO")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("NOTI")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("NOTICE")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("WARN")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("WARNING")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("ERRO")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("ERROR")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("CRIT")
		t.Assert(enablePrintRunLog(), false)
		_ = internal.SetLevelStr("CRITICAL")
		t.Assert(enablePrintRunLog(), false)
	})
}
