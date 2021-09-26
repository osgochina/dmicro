# EventBus(事件总线)

## 概念

事件模式是一种经过了充分测试的可靠机制，是一种非常适用于解耦的机制，分别存在以下 3 种角色：

1. 事件(Event) 是传递于应用代码与 监听器(Listener) 之间的通讯对象
2. 监听器(Listener) 是用于监听 事件(Event) 的发生的监听对象
3. 事件调度器(Manager) 用于触发事件(Event),管理 监听器(Listener) 与 事件(Event) 之间的关系的管理者对象

用通俗易懂的例子来说明就是，假设我们存在一个 `register()` 方法用于注册一个账号，
在账号注册成功后我们可以通过事件调度器触发 `UserRegistered` 事件， 由监听器监听该事件的发生，在触发时进行某些操作，
比如发送用户注册成功短信，在业务发展的同时我们可能会希望在用户注册成功之后做更多的事情，
比如发送用户注册成功的邮件等待，此时我们就可以通过再增加一个监听器监听 `UserRegistered` 事件即可，
无需在 `register()` 方法内部增加与之无关的代码。


## 功能特性

* 支持自定义定义事件对象
* 支持对一个事件添加多个监听器
* 支持设置事件监听的优先级，优先级越高越先触发
* 支持根据事件名称前缀 `PREFIX.*` 来进行一组事件监听.
    > 注册app.* 事件的监听，触发 app.run app.end 时，都将同时会触发 app.* 事件
* 支持使用通配符 * 来监听全部事件的触发

## 快速使用

```go
package main

import (
  "fmt"
  "github.com/osgochina/dmicro/eventbus"
  "github.com/osgochina/dmicro/logger"
)

func main() {
  var event1 eventbus.IEvent
  // 注册监听器，监听事件event1
  err := eventbus.Listen("event1",eventbus.ListenerFunc(func(e eventbus.IEvent) error {
    fmt.Printf("listen1 process event: %s \n", e.Name())
    return nil
  }),eventbus.Normal)
  if err!=nil {
    logger.Fatal(err)
  }

  err = eventbus.Listen("event1",eventbus.ListenerFunc(func(e eventbus.IEvent) error {
    fmt.Printf("listen2 process event: %s \n", e.Name())
    return nil
  }),eventbus.High)
  if err!=nil {
    logger.Fatal(err)
  }
  // 执行事件的时候，会优先执行listen2，因为它的优先度高
  event1,err = eventbus.Fire("event1", map[interface{}]interface{}{"arg0":"val0","arg1":"val1"})
  if err!=nil {
    logger.Fatal(err)
  }
  // 输出 val0
  fmt.Println(event1.Get("arg0"))
}

```

## 使用匿名函数监听

```go
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
	err := eventbus.Listen("event1",eventbus.ListenerFunc(anonymity),eventbus.Normal)
	if err!=nil {
		logger.Fatal(err)
	}
	_,err = eventbus.Fire("event1", nil)
	if err!=nil {
		logger.Fatal(err)
	}
}

```

## 使用结构体方法监听

```go
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
	e.Set("result","ok")
	return nil
}


// 使用结构体方法作为监听器
func main() {
	var event1 eventbus.IEvent
	err := eventbus.Listen("event1",&MyListener{},eventbus.Normal)
	if err!=nil {
		logger.Fatal(err)
	}
	event1,err = eventbus.Fire("event1", nil)
	if err!=nil {
		logger.Fatal(err)
	}
	// 输出 ok
	fmt.Println(event1.Get("result"))
}
```

## 同时监听多个事件

```go
package main

import (
	"fmt"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
)

// MyListenerMulti 需要实现接口 eventbus.IListener
type MyListenerMulti struct {}

func (that *MyListenerMulti)Process(e eventbus.IEvent) error {
	fmt.Println(e.Name())
	e.Set("result","ok")
	return nil
}

func (that *MyListenerMulti)Listen() []interface{} {
	return []interface{}{
		"event1",
		"event2",
	}
}

// 通过结构体注册多个事件到同一个监听器
func main() {
	err := eventbus.Subscribe(&MyListenerMulti{},eventbus.High)
	if err!=nil {
		logger.Fatal(err)
	}
	errs := eventbus.PublishBatch("event1","event2")
	if len(errs) > 0 {
		logger.Fatal(errs)
	}
}
```

## 自定义事件监听

```go
package main

import (
	"fmt"
	"github.com/osgochina/dmicro/eventbus"
	"github.com/osgochina/dmicro/logger"
)

// MyListenerCustom 需要实现接口 eventbus.IListener
type MyListenerCustom struct {}

func (that *MyListenerCustom)Process(e eventbus.IEvent) error {
	if m, ok := e.(*MyEvent); ok {
		fmt.Println(e.Name(),m.CustomData())
	}else{
		fmt.Println(e.Name())
	}
	e.Set("result","ok")
	return nil
}

func (that *MyListenerCustom)Listen() []interface{} {
	return []interface{}{
		"event1",
		"event2",
	}
}
type MyEvent struct {
	eventbus.Event
	customData string
}

func (that *MyEvent)CustomData() string {
	return that.customData
}


func main() {
    // 通过结构体注册多个事件到同一个监听器
	err := eventbus.Subscribe(&MyListenerCustom{},eventbus.High)
	if err!=nil {
		logger.Fatal(err)
	}
	customEvent:=&MyEvent{customData: "clownfish"}
	customEvent.SetName("event2")
	err = eventbus.AddEvent(customEvent)
	if err!=nil {
		logger.Fatal(err)
	}
	errs := eventbus.PublishBatch("event2","event1")
	if len(errs) > 0 {
		logger.Fatal(errs)
	}
}
```

