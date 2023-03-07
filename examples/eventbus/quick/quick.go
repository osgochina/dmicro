package main

import (
	"context"
	"fmt"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
)

func main() {
	var event1 eventbus.IEvent
	// 注册监听器，监听事件event1
	err := eventbus.Listen("event1", eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		fmt.Printf("listen1 process event: %s \n", e.Name())
		return nil
	}), eventbus.Normal)
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}

	err = eventbus.Listen("event1", eventbus.ListenerFunc(func(e eventbus.IEvent) error {
		fmt.Printf("listen2 process event: %s \n", e.Name())
		return nil
	}), eventbus.High)
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}
	// 执行事件的时候，会优先执行listen2，因为它的优先度高
	event1, err = eventbus.Fire("event1", map[interface{}]interface{}{"arg0": "val0", "arg1": "val1"})
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}
	// 输出 val0
	fmt.Println(event1.Get("arg0"))
}
