package main

import (
	"fmt"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
)

// MyListener 需要实现接口 eventbus.BaseListener
type MyListener struct {
	name string
}

func (that *MyListener) Process(e eventbus.IEvent) error {
	that.name = e.Name()
	e.Set("result", "ok")
	return nil
}

// 使用结构体方法作为监听器
func main() {
	var event1 eventbus.IEvent
	err := eventbus.Listen("event1", &MyListener{}, eventbus.Normal)
	if err != nil {
		logger.Fatal(err)
	}
	event1, err = eventbus.Fire("event1", nil)
	if err != nil {
		logger.Fatal(err)
	}
	// 输出 ok
	fmt.Println(event1.Get("result"))
}
