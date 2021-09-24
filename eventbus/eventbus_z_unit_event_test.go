package eventbus

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func TestEvent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		e1 := NewEvent("test1", nil)
		e1.SetName("test1")
		e1.SetData(map[interface{}]interface{}{
			"arg0": "val0",
		})

		e1.Set("arg1", "val1")
		t.Assert(e1.IsAborted(), false)
		e1.Abort(true)
		t.Assert(e1.IsAborted(), true)
		t.Assert(e1.Name(), "test1")
		t.Assert(e1.Get("arg1"), "val1")
		t.Assert(e1.Get("not_exist"), nil)
		e1.Set("arg1", "new val1")
		t.Assert(e1.Get("arg1"), "new val1")

		e2 := &Event{}
		e2.Set("k", "v")
		e2.Set("k1", "v1")
		t.Assert(e2.Get("k"), "v")
		var contains = false
		for k, _ := range e2.Data() {
			if k == "k1" {
				contains = true
			}
		}
		t.Assert(contains, true)
	})
}

func TestAddEvent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := New("default")
		defer manager.Clear()

		err := manager.AddEvent(&Event{})
		t.AssertNE(err, nil)

		_, ok := manager.GetEvent("evt1")
		t.Assert(ok, false)

		e := NewEvent("evt1", map[interface{}]interface{}{"k1": "val1"})
		err = manager.AddEvent(e)
		t.Assert(err, nil)

		err = NewEvent("evt2", nil).AttachTo(manager)
		t.Assert(err, nil)

		t.Assert(e.IsAborted(), false)

		t.Assert(manager.HasEvent("evt1"), true)
		t.Assert(manager.HasEvent("evt2"), true)
		t.Assert(manager.HasEvent("not-exist"), false)

		r1, ok := manager.GetEvent("evt1")
		t.Assert(ok, true)
		t.Assert(r1, e)

		manager.RemoveEvent("evt2")
		t.Assert(manager.HasEvent("evt2"), false)
		manager.RemoveEvents()
		t.Assert(manager.HasEvent("evt1"), false)

	})
}
