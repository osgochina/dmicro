package eventbus

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/v2/test/gtest"
	"sync"
	"testing"
	"time"
)

type testListener struct {
	listen []interface{}
	data   string
}

func (that *testListener) Process(e IEvent) error {
	if ret := e.Get("result"); ret != nil {
		str := ret.(string) + fmt.Sprintf(" -> %s(%s)", e.Name(), that.data)
		e.Set("result", str)
	} else {
		e.Set("result", fmt.Sprintf("process: %s(%s)", e.Name(), that.data))
	}
	return nil
}

func (that *testListener) Listen() []interface{} {
	return that.listen
}

func TestManagerListen(t *testing.T) {
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

func TestManagerFire(t *testing.T) {
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
		}), AboveNormal)
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

func TestManagerPublish(t *testing.T) {
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
		time.Sleep(time.Second * 1)
		t.Assert(buf, "event: evt1, params: n=inhere")
	})
}

func TestManagerPublish2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := New("default")
		defer manager.Clear()

		e1 := NewEvent("e1", nil)
		err := manager.AddEvent(e1)
		t.Assert(err, nil)

		err = manager.Listen("e1", &testListener{data: "Welcome"}, Min)
		t.Assert(err, nil)
		err = manager.Listen("e1", &testListener{data: "Hello"}, High)
		t.Assert(err, nil)
		err = manager.Subscribe(&testListener{data: "ClownFish", listen: []interface{}{"e1"}}, BelowNormal)
		t.Assert(err, nil)

		err = manager.Publish(e1)
		t.Assert(err, nil)
		t.Assert(e1.Get("result"), "process: e1(Hello) -> e1(ClownFish) -> e1(Welcome)")
		e1.SetName("e2")
		err = manager.Publish(e1)
		t.Assert(err, nil)
		manager.Clear()
	})
}

func TestManagerFireWithWildcard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := New("default")
		defer manager.Clear()
		buf := new(bytes.Buffer)
		var eventName = "site.clownfish.www"
		process := ListenerFunc(func(e IEvent) error {
			_, _ = fmt.Fprintf(buf, "%s-%s|", e.Name(), e.Get("user"))
			return nil
		})
		err := manager.Listen("site.clownfish.*", process)
		t.Assert(err, nil)
		err = manager.Listen(eventName, process)
		t.Assert(err, nil)

		_, err = manager.Fire(eventName, map[interface{}]interface{}{"user": "liuzhiming"})
		t.Assert(err, nil)
		t.Assert(buf, "site.clownfish.www-liuzhiming|site.clownfish.www-liuzhiming|")

		buf.Reset()
		err = manager.Listen("*", process)
		t.Assert(err, nil)

		err = manager.Publish(NewEvent(eventName, map[interface{}]interface{}{"user": "lzm"}))
		t.Assert(err, nil)
		t.Assert(buf, "site.clownfish.www-lzm|site.clownfish.www-lzm|site.clownfish.www-lzm|")
	})
}

func TestListenGroupEvent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := New("default")
		defer manager.Clear()

		e1 := NewEvent("app.evt1", map[interface{}]interface{}{"buf": new(bytes.Buffer)})
		err := e1.AttachTo(manager)
		t.Assert(err, nil)
		l2 := ListenerFunc(func(e IEvent) error {
			e.Get("buf").(*bytes.Buffer).WriteString(" > 2 " + e.Name())
			return nil
		})
		l3 := ListenerFunc(func(e IEvent) error {
			e.Get("buf").(*bytes.Buffer).WriteString(" > 3 " + e.Name())
			return nil
		})
		err = manager.Listen("app.evt1", ListenerFunc(func(e IEvent) error {
			e.Get("buf").(*bytes.Buffer).WriteString("Hi > 1 " + e.Name())
			return nil
		}))
		t.Assert(err, nil)
		err = manager.Listen("app.*", l2)
		t.Assert(err, nil)
		err = manager.Listen("*", l3)
		t.Assert(err, nil)

		buf := e1.Get("buf").(*bytes.Buffer)
		e, err := manager.Fire("app.evt1", nil)
		t.Assert(err, nil)
		t.Assert(e1, e)
		t.Assert(buf, "Hi > 1 app.evt1 > 2 app.evt1 > 3 app.evt1")

		manager.RemoveListenersByName("app.*")
		t.Assert(len(manager.ListenedNames()), 2)
		err = manager.Listen("app.*", ListenerFunc(func(e IEvent) error {
			return fmt.Errorf("an error")
		}))
		t.Assert(err, nil)
		buf.Reset()

		e, err = manager.Fire("app.evt1", nil)
		t.AssertNE(err, nil)
		t.Assert(buf.String(), "Hi > 1 app.evt1")

		buf.Reset()
		manager.RemoveListenersByName("app.*")
		e, err = manager.Fire("app.evt1", nil)
		t.Assert(err, nil)
		t.Assert(buf.String(), "Hi > 1 app.evt1 > 3 app.evt1")

		manager.RemoveListeners(l3)
		err = manager.Listen("app.*", l2)
		t.Assert(err, nil)
		err = manager.Listen("*", ListenerFunc(func(e IEvent) error {
			return fmt.Errorf("an error")
		}))
		t.Assert(err, nil)
		t.Assert(len(manager.ListenedNames()), 3)

		buf.Reset()

		e, err = manager.Fire("app.evt1", nil)
		t.AssertNE(err, nil)
		t.Assert(e1, e)
		t.Assert(buf.String(), "Hi > 1 app.evt1 > 2 app.evt1")
		buf.Reset()
	})
}

func TestManagerAsyncFire(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		manager := New("default")
		defer manager.Clear()

		err := manager.Listen("e1", ListenerFunc(func(e IEvent) error {
			t.Assert(map[string]interface{}{"k": "v"}, e.Data())
			e.Set("nk", "nv")
			return nil
		}))
		t.Assert(err, nil)

		e1 := NewEvent("e1", map[interface{}]interface{}{"k": "v"})
		manager.AsyncPublish(e1)
		time.Sleep(time.Second / 10)

		t.Assert(e1.Get("nk"), "nv")
		var wg sync.WaitGroup
		err = manager.Listen("e2", ListenerFunc(func(e IEvent) error {
			defer wg.Done()
			t.Assert(e.Get("k"), "v")
			return nil
		}))
		t.Assert(err, nil)

		wg.Add(1)
		e1.SetName("e2")
		manager.AsyncPublish(e1)
		wg.Wait()
	})

}
