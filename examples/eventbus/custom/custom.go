package main

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
)

// MyListenerCustom 需要实现接口 eventbus.IListener
type MyListenerCustom struct{}

func (that *MyListenerCustom) Process(e eventbus.IEvent) error {
	if m, ok := e.(*MyEvent); ok {
		fmt.Println(e.Name(), m.CustomData())
	} else {
		fmt.Println(e.Name())
	}
	e.Set("result", "ok")
	return nil
}

func (that *MyListenerCustom) Listen() []interface{} {
	return []interface{}{
		"event1",
		"event2",
	}
}

type MyEvent struct {
	eventbus.Event
	customData string
}

func (that *MyEvent) CustomData() string {
	return that.customData
}

func main() {
	// 通过结构体注册多个事件到同一个监听器
	err := eventbus.Subscribe(&MyListenerCustom{}, eventbus.High)
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}
	customEvent := &MyEvent{customData: "clownfish"}
	customEvent.SetName("event2")
	err = eventbus.AddEvent(customEvent)
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}
	errs := eventbus.PublishBatch("event2", "event1")
	if len(errs) > 0 {
		logger.Fatal(context.TODO(), errs)
	}
}
