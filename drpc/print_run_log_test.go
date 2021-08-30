package drpc

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/logger"
	"testing"
)

func TestEnablePrintRunLog(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_ = logger.SetLevelStr("ALL")
		t.Assert(enablePrintRunLog(), true)
		_ = logger.SetLevelStr("DEV")
		t.Assert(enablePrintRunLog(), true)
		_ = logger.SetLevelStr("DEVELOP")
		t.Assert(enablePrintRunLog(), true)
		_ = logger.SetLevelStr("PROD")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("PRODUCT")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("DEBU")
		t.Assert(enablePrintRunLog(), true)
		_ = logger.SetLevelStr("DEBUG")
		t.Assert(enablePrintRunLog(), true)
		_ = logger.SetLevelStr("INFO")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("NOTI")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("NOTICE")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("WARN")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("WARNING")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("ERRO")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("ERROR")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("CRIT")
		t.Assert(enablePrintRunLog(), false)
		_ = logger.SetLevelStr("CRITICAL")
		t.Assert(enablePrintRunLog(), false)
	})
}
