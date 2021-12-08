package eventbus

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

var emptyListener = func(e IEvent) error {
	return nil
}

var emptyListener2 = func(e IEvent) error {
	return nil
}

func TestListenerQueueRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		obj := newListenerQueue()
		obj.Add(NewListenerItem(ListenerFunc(emptyListener2), 1))
		lt := ListenerFunc(emptyListener)
		for n := 0; n < 10; n++ {
			//name := fmt.Sprintf("evt%d", n)
			obj.Add(NewListenerItem(lt, n))
		}
		t.Assert(obj.Len(), 11)
		obj.Remove(lt)
		t.Assert(obj.Len(), 1)
	})
}

var testBuf = new(bytes.Buffer)

type testListenerList struct{}

var testevt3 = NewEvent("evt3", nil)

func (that *testListenerList) Listen() []interface{} {
	return []interface{}{
		"evt1",
		"evt2",
		testevt3,
	}
}

func (that *testListenerList) Process(event IEvent) error {
	event.Set("name", event.Name())
	_, _ = fmt.Fprintf(testBuf, "%s|", event.Name())
	return nil
}

func TestAddSubscriber(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		e := &testListenerList{}
		err := Subscribe(e, Low)
		t.Assert(err, nil)
		t.Assert(HasListeners("evt1"), true)
		t.Assert(HasListeners("evt2"), true)
		t.Assert(HasListeners("evt3"), true)
		t.Assert(HasListeners("evt4"), false)

		evt1, err := Fire("evt1", nil)
		t.Assert(err, nil)
		t.Assert(evt1.Get("name"), "evt1")
		t.Assert(testBuf.String(), "evt1|")
		evt2, err := Fire("evt2", nil)
		t.Assert(err, nil)
		t.Assert(evt2.Get("name"), "evt2")
		t.Assert(testBuf.String(), "evt1|evt2|")
		evt3, err := Fire("evt3", nil)
		t.Assert(err, nil)
		t.Assert(evt3.Get("name"), "evt3")
		t.Assert(testBuf.String(), "evt1|evt2|evt3|")
		t.Assert(evt3, testevt3)
		testBuf.Reset()
		errs := PublishBatch("evt1", "evt2", testevt3)
		t.Assert(len(errs), 0)
		t.Assert(testBuf.String(), "evt1|evt2|evt3|")

		testBuf.Reset()
		RemoveListenersByName("evt2")
		errs = PublishBatch("evt1", "evt2", testevt3)
		t.Assert(len(errs), 0)
		t.Assert(testBuf.String(), "evt1|evt3|")
		testBuf.Reset()

		err = Publish(testevt3)
		t.Assert(err, nil)
		t.Assert(testBuf.String(), "evt3|")

	})
}
