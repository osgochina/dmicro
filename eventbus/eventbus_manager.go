package eventbus

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"strings"
)

// EventBus 事件管理器
type EventBus struct {
	//管理器名称
	name string
	// 已注册的事件列表
	events *gmap.StrAnyMap
	//已监听的监听列表
	listeners *gmap.StrAnyMap
}

// New 创建事件管理器对象
func New(name string) *EventBus {
	em := &EventBus{
		name:      name,
		events:    gmap.NewStrAnyMap(true),
		listeners: gmap.NewStrAnyMap(true),
	}
	return em
}

//================================================== listener ========================================================//

// Subscribe 注册监听事件
func (that *EventBus) Subscribe(listener IListener, priority ...int) (err error) {
	for _, val := range listener.Listen() {
		switch val.(type) {
		case string:
			err = that.Listen(val.(string), listener, priority...)
			if err != nil {
				return err
			}
			break
		case IEvent:
			name := val.(IEvent).Name()
			//判断事件是否存在，如果不存在，则添加该事件
			if !that.HasEvent(name) {
				err = that.AddEvent(val.(IEvent))
				if err != nil {
					return err
				}
			}
			err = that.Listen(name, listener, priority...)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// UnSubscribe 取消监听事件
func (that *EventBus) UnSubscribe(listener IListener) {
	for _, val := range listener.Listen() {
		switch val.(type) {
		case string:
			that.UnListen(val.(string), listener)
			break
		case IEvent:
			that.UnListen(val.(IEvent).Name(), listener)
			break
		}
	}
	return
}

// Listen 监听事件
func (that *EventBus) Listen(name string, listener BaseListener, priority ...int) (err error) {
	pv := Normal
	if len(priority) > 0 {
		pv = priority[0]
	}
	if name != Wildcard {
		name, err = checkName(name)
		if err != nil {
			return err
		}
	}

	if listener == nil {
		return gerror.New("event: the event '" + name + "' listener cannot be empty")
	}
	listen, found := that.listeners.Search(name)
	if found {
		listen.(*ListenerQueue).Add(NewListenerItem(listener, pv))
	} else {
		obj := newListenerQueue()
		obj.Add(NewListenerItem(listener, pv))
		that.listeners.Set(name, obj)
	}
	return nil
}

// UnListen 移除事件监听
func (that *EventBus) UnListen(name string, listener BaseListener) {
	if val, ok := that.listeners.Search(name); ok {
		lq := val.(*ListenerQueue)
		lq.Remove(listener)

		// 如果从监听列表移除监听事件后，监听列表为空，则把该事件名全部移除
		if lq.IsEmpty() {
			that.listeners.Remove(name)
		}
	}
	return
}

// RemoveListenersByName 清除所有指定事件名的事件监听器
func (that *EventBus) RemoveListenersByName(name string) {
	l, ok := that.listeners.Search(name)
	if ok {
		l.(*ListenerQueue).Clear()
		that.listeners.Remove(name)
	}
}

// RemoveListeners 移除监听方法
func (that *EventBus) RemoveListeners(listener BaseListener) {
	that.listeners.LockFunc(func(m map[string]interface{}) {
		for name, v := range m {
			v.(*ListenerQueue).Remove(listener)
			if v.(*ListenerQueue).IsEmpty() {
				delete(m, name)
			}
		}
	})
}

//================================================== Listener ========================================================//

// HasListeners 判断是否存在指定名称的监听
func (that *EventBus) HasListeners(name string) bool {
	return that.listeners.Contains(name)
}

// Listeners 获取监听列表
func (that *EventBus) Listeners() map[string]*ListenerQueue {
	var result = make(map[string]*ListenerQueue, that.listeners.Size())
	for k, v := range that.listeners.Map() {
		result[k] = v.(*ListenerQueue)
	}
	return result
}

// ListenersByName 获取指定事件名称的监听列表
func (that *EventBus) ListenersByName(name string) *ListenerQueue {
	result := that.listeners.Get(name)
	if result == nil {
		return nil
	}
	return result.(*ListenerQueue)
}

// ListenersCount 获取指定名称的监听列表数量
func (that *EventBus) ListenersCount(name string) int {
	result := that.listeners.Get(name)
	if result != nil {
		return result.(*ListenerQueue).Len()
	}
	return 0
}

// ListenedNames 获取监听的事件名列表
func (that *EventBus) ListenedNames() []string {
	return that.listeners.Keys()
}

//================================================== event ===========================================================//

// AddEvent 添加自定义事件
func (that *EventBus) AddEvent(e IEvent) error {
	name, err := checkName(e.Name())
	if err != nil {
		return err
	}
	if that.events.Contains(name) {
		return gerror.Newf("event %s is exist!", name)
	}
	that.events.Set(name, e)
	return nil
}

// GetEvent 根据name获取自定义事件
func (that *EventBus) GetEvent(name string) (IEvent, bool) {
	v, ok := that.events.Search(name)
	if ok {
		return v.(IEvent), ok
	}
	return nil, ok
}

// HasEvent 判断自定义事件是否存在
func (that *EventBus) HasEvent(name string) bool {
	return that.events.Contains(name)
}

// RemoveEvent 移除自定义事件
func (that *EventBus) RemoveEvent(name string) {
	that.events.Remove(name)
}

// RemoveEvents 移除所有自定义事件
func (that *EventBus) RemoveEvents() {
	that.events.Clear()
}

//================================================== Publish =========================================================//

// Fire 使用事件名触发事件
// name: 事件名
// params:需要传递给事件的参数
func (that *EventBus) Fire(name string, params map[interface{}]interface{}) (e IEvent, err error) {
	name, err = checkName(name)
	if err != nil {
		return nil, err
	}

	// 判断要触发的事件名是否存在,如果存在监听"*"所有事件的触发器。则继续执行
	if !that.HasListeners(name) && !that.HasListeners(Wildcard) {
		return nil, nil
	}

	// 判断要触发的事件是否存在
	if ev, ok := that.events.Search(name); ok {
		e = ev.(IEvent)
		if params != nil {
			e.SetData(params)
		}
		err = that.Publish(e)
		return e, err
	}
	// 创建一个事件对象，并且触发它
	e = NewEvent(name, params)
	err = that.Publish(e)
	return e, err
}

// Publish 触发事件
func (that *EventBus) Publish(e IEvent) error {
	// 把中断标记设置为false
	e.Abort(false)
	name := e.Name()
	//通过事件名称，查找监听的方法，触发执行
	queueListeners, found := that.listeners.Search(name)
	if found && queueListeners != nil && queueListeners.(*ListenerQueue).Len() > 0 {
		lq := queueListeners.(*ListenerQueue)
		for _, item := range lq.Items() {
			err := item.Listener.Process(e)
			if err != nil || e.IsAborted() {
				return err
			}
		}
	}
	// 查找分组监听的情况，比如"app.*" "app.cache.*"
	// 比如："app.run"事件会触发"app.*"的监听
	pos := strings.LastIndexByte(name, '.')
	if pos > 0 && pos < len(name) {
		groupName := name[:pos+1] + Wildcard // "app.*"
		queueListeners, found = that.listeners.Search(groupName)
		if found && queueListeners != nil && queueListeners.(*ListenerQueue).Len() > 0 {
			lq := queueListeners.(*ListenerQueue)
			for _, item := range lq.Items() {
				err := item.Listener.Process(e)
				if err != nil || e.IsAborted() {
					return err
				}
			}
		}
	}
	// 获取队列的完全匹配
	queueListeners, found = that.listeners.Search(Wildcard)
	if found && queueListeners != nil && queueListeners.(*ListenerQueue).Len() > 0 {
		lq := queueListeners.(*ListenerQueue)
		for _, item := range lq.Items() {
			err := item.Listener.Process(e)
			if err != nil || e.IsAborted() {
				return err
			}
		}
	}
	return nil
}

// PublishBatch 批量触发事件
// Usage:
// 	PublishBatch("name1", "name2", &MyEvent{})
func (that *EventBus) PublishBatch(es ...interface{}) (ers []error) {
	var err error
	for _, e := range es {
		if name, ok := e.(string); ok {
			_, err = that.Fire(name, nil)
		} else if evt, ok := e.(IEvent); ok {
			err = that.Publish(evt)
		}
		if err != nil {
			ers = append(ers, err)
		}
	}
	return
}

// AsyncPublish 异步触发事件
func (that *EventBus) AsyncPublish(e IEvent) {
	go func(e IEvent) {
		_ = that.Publish(e)
	}(e)
}

// Clear 清空事件管理对象
func (that *EventBus) Clear() {
	that.Reset()
}

// Reset 重置事件管理对象
func (that *EventBus) Reset() {
	that.name = ""
	that.listeners.Clear()
	that.events.Clear()
}
