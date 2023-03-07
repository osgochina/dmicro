package main

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
)

// MyListenerMulti 需要实现接口 eventbus.IListener
type MyListenerMulti struct{}

func (that *MyListenerMulti) Process(e eventbus.IEvent) error {
	fmt.Println(e.Name())
	e.Set("result", "ok")
	return nil
}

func (that *MyListenerMulti) Listen() []interface{} {
	return []interface{}{
		"event1",
		"event2",
	}
}

func main() {
	// 通过结构体注册多个事件到同一个监听器
	err := eventbus.Subscribe(&MyListenerMulti{}, eventbus.High)
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}
	errs := eventbus.PublishBatch("event1", "event2")
	if len(errs) > 0 {
		logger.Fatal(context.TODO(), errs)
	}
}
