package eventbus

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/event"
	"testing"
	"time"
)

func TestManager_Listen(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := New("default")
		defer manager.Clear()

		err := manager.Listen("", ListenerFunc(emptyListener), 0)
		t.AssertNE(err, nil)
		err = manager.Listen("name", nil, 0)
		t.AssertNE(err, nil)
		err = manager.Listen("++addf", ListenerFunc(emptyListener), 0)
		t.AssertNE(err, nil)

		err = manager.Listen("n1", ListenerFunc(emptyListener), 0)
		t.Assert(err, nil)
	})
}

func TestManager_Fire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := New("default")
		defer manager.Clear()
		buf := new(bytes.Buffer)
		fn := func(e IEvent) error {
			_, _ = fmt.Fprintf(buf, "event: %s", e.Name())
			return nil
		}
		err := manager.Listen("evt1", ListenerFunc(fn), 0)
		t.Assert(err, nil)
		err = manager.Listen("evt1", ListenerFunc(emptyListener), High)
		t.Assert(err, nil)
		e, err := manager.Fire("evt1", nil)
		t.Assert(err, nil)
		t.Assert(e.Name(), "evt1")
		t.Assert(buf.String(), "event: evt1")

		err = NewEvent("evt2", nil).AttachTo(manager)
		t.Assert(err, nil)

		err = manager.Listen("evt2", ListenerFunc(func(e IEvent) error {
			t.Assert(e.Name(), "evt2")
			t.Assert(e.Get("k"), "v")
			return nil
		}), event.AboveNormal)
		t.Assert(err, nil)

		e, err = manager.Fire("evt2", map[interface{}]interface{}{"k": "v"})
		t.Assert(err, nil)

		t.Assert(e.Name(), "evt2")
		t.Assert(e.Data(), map[interface{}]interface{}{"k": "v"})

		manager.Reset()

		e, err = manager.Fire("not-exist", nil)
		t.Assert(err, nil)
		t.Assert(e, nil)
	})
}

func TestManager_Publish(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := New("default")
		defer manager.Clear()
		buf := new(bytes.Buffer)

		evt1 := NewEvent("evt1", nil).SetData(map[interface{}]interface{}{"n": "inhere"})
		err := manager.AddEvent(evt1)
		t.Assert(err, nil)

		t.Assert(manager.HasEvent("evt1"), true)
		t.Assert(manager.HasEvent("not-exist"), false)

		err = manager.Listen("evt1", ListenerFunc(func(e IEvent) error {
			_, _ = fmt.Fprintf(buf, "event: %s, params: n=%s", e.Name(), e.Get("n"))
			return nil
		}), Normal)
		t.Assert(err, nil)

		t.Assert(manager.HasEvent("evt1"), true)
		t.Assert(manager.HasEvent("not-exist"), false)

		err = manager.Publish(evt1)
		t.Assert(err, nil)
		t.Assert(buf, "event: evt1, params: n=inhere")

		buf.Reset()

		manager.AsyncPublish(evt1)
		time.Sleep(time.Second * 10)
		t.Assert(buf, "event: evt1, params: n=inhere")
	})
}
