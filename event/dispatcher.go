package event

import "strings"

// Dispatcher 通过事件名触发监听
func (that *Manager) Dispatcher(name string, params map[interface{}]interface{}) (e IEvent, err error) {
	name, err = checkName(name)
	if err != nil {
		return nil, err
	}

	// 判断要触发的事件名是否存在,如果存在监听"*"所有事件的触发器。则继续执行
	if that.HasListeners(name) == false && that.HasListeners(Wildcard) == false {
		return nil, nil
	}

	// 判断要触发的事件是否存在
	if e, ok := that.events[name]; ok {
		if params != nil {
			e.SetData(params)
		}

		err = that.DispatcherByEvent(e)
		return e, err
	}
	// 创建一个事件对象，并且触发它
	e = NewEvent(name, params)
	err = that.DispatcherByEvent(e)
	return e, err
}

// DispatcherByEvent 触发事件
func (that *Manager) DispatcherByEvent(e IEvent) error {
	e.Abort(false)
	name := e.Name()

	//通过事件名称，查找监听的方法，触发执行
	queueListeners := that.listeners.GetListenersForEvent(name)
	if queueListeners != nil && queueListeners.Len() > 0 {
		for i := 0; i < queueListeners.Len(); i++ {
			item := queueListeners.Pop()
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
		queueListeners = that.listeners.GetListenersForEvent(groupName)
		if queueListeners != nil && queueListeners.Len() > 0 {
			for i := 0; i < queueListeners.Len(); i++ {
				item := queueListeners.Pop()
				err := item.Listener.Process(e)
				if err != nil || e.IsAborted() {
					return err
				}
			}
		}
	}

	// 获取队列的完全匹配
	queueListeners = that.listeners.GetListenersForEvent(Wildcard)
	if queueListeners != nil && queueListeners.Len() > 0 {
		for i := 0; i < queueListeners.Len(); i++ {
			item := queueListeners.Pop()
			err := item.Listener.Process(e)
			if err != nil || e.IsAborted() {
				return err
			}
		}
	}
	return nil
}

// AsyncDispatcher 异步触发事件
func (that *Manager) AsyncDispatcher(e IEvent) {
	go func(e IEvent) {
		_ = that.DispatcherByEvent(e)
	}(e)
}
