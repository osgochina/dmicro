package event_test

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/event"
	"testing"
)

var emptyListener = func(e event.IEvent) error {
	return nil
}

func TestManager_On(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := event.NewManager("default")
		defer manager.Clear()

		err := manager.On("", event.ListenerFunc(emptyListener), 0)
		t.AssertNE(err, nil)
		err = manager.On("name", nil, 0)
		t.AssertNE(err, nil)
		err = manager.On("++addf", event.ListenerFunc(emptyListener), 0)
		t.AssertNE(err, nil)

		err = manager.On("n1", event.ListenerFunc(emptyListener), 0)
		t.Assert(err, nil)
	})
}
