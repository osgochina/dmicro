package main

import (
	"fmt"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
)

var anonymity = func(e eventbus.IEvent) error {
	fmt.Printf("process event: %s \n", e.Name())
	return nil
}

// 使用匿名函数注册监听器
func main() {
	err := eventbus.Listen("event1", eventbus.ListenerFunc(anonymity), eventbus.Normal)
	if err != nil {
		logger.Fatal(err)
	}
	_, err = eventbus.Fire("event1", nil)
	if err != nil {
		logger.Fatal(err)
	}
}
