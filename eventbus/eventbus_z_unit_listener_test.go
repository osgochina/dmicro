package eventbus

import (
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

func TestListenerQueue_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		obj := newListenerQueue()
		obj.Push(NewListenerItem("evt", ListenerFunc(emptyListener2), 1))
		lt := ListenerFunc(emptyListener)
		for n := 0; n < 10; n++ {
			name := fmt.Sprintf("evt%d", n)
			obj.Push(NewListenerItem(name, lt, n))
		}
		t.Assert(obj.Len(), 11)
		obj.Remove(lt)
		t.Assert(obj.Len(), 1)
	})
}
