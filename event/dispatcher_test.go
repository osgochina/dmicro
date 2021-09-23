package event_test

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/event"
	"testing"
	"time"
)

func TestManager_Dispatcher(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := event.NewManager("default")
		defer manager.Clear()
		buf := new(bytes.Buffer)
		fn := func(e event.IEvent) error {
			_, _ = fmt.Fprintf(buf, "event: %s", e.Name())
			return nil
		}
		err := manager.On("evt1", event.ListenerFunc(fn), 0)
		t.Assert(err, nil)
		err = manager.On("evt1", event.ListenerFunc(emptyListener), event.High)
		t.Assert(err, nil)
		e, err := manager.Dispatcher("evt1", nil)
		t.Assert(err, nil)
		t.Assert(e.Name(), "evt1")
		t.Assert(buf.String(), "event: evt1")

		err = event.NewEvent("evt2", nil).AttachTo(manager)
		t.Assert(err, nil)

		err = manager.On("evt2", event.ListenerFunc(func(e event.IEvent) error {
			t.Assert(e.Name(), "evt2")
			t.Assert(e.Get("k"), "v")
			return nil
		}), event.AboveNormal)
		t.Assert(err, nil)

		e, err = manager.Dispatcher("evt2", map[interface{}]interface{}{"k": "v"})
		t.Assert(err, nil)

		t.Assert(e.Name(), "evt2")
		t.Assert(e.Data(), map[interface{}]interface{}{"k": "v"})

		manager.Reset()

		e, err = manager.Dispatcher("not-exist", nil)
		t.Assert(err, nil)
		t.Assert(e, nil)
	})
}

func TestManager_DispatcherByEvent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := event.NewManager("default")
		defer manager.Clear()
		buf := new(bytes.Buffer)

		evt1 := event.NewEvent("evt1", nil).Fill(nil, map[interface{}]interface{}{"n": "inhere"})
		err := manager.AddEvent(evt1)
		t.Assert(err, nil)

		t.Assert(manager.HasEvent("evt1"), true)
		t.Assert(manager.HasEvent("not-exist"), false)

		err = manager.Listen("evt1", event.ListenerFunc(func(e event.IEvent) error {
			_, _ = fmt.Fprintf(buf, "event: %s, params: n=%s", e.Name(), e.Get("n"))
			return nil
		}), event.Normal)
		t.Assert(err, nil)

		t.Assert(manager.HasEvent("evt1"), true)
		t.Assert(manager.HasEvent("not-exist"), false)

		err = manager.DispatcherByEvent(evt1)
		t.Assert(err, nil)
		t.Assert(buf, "event: evt1, params: n=inhere")

		buf.Reset()

		manager.AsyncDispatcher(evt1)
		time.Sleep(time.Second)
		t.Assert(buf, "event: evt1, params: n=inhere")
	})
}
