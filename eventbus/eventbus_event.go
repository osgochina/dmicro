package eventbus

import "github.com/gogf/gf/container/gmap"

// IEvent 事件的接口
type IEvent interface {
	// Name 获取事件名
	Name() string
	// Get 获取事件携带的参数
	Get(key interface{}) interface{}
	// Set 设置事件参数
	Set(key interface{}, val interface{})
	// Data 获取事件的全部参数
	Data() map[interface{}]interface{}
	// SetData 批量设置事件参数
	SetData(map[interface{}]interface{}) IEvent
	// Abort 设置是否中途停止，设置为true，则表示执行到该监听器为止，后续监听器不在执行
	Abort(bool)
	// IsAborted 判断是否要继续触发后续监听器
	IsAborted() bool
}

// Event 事件的基础实现
type Event struct {
	name    string
	data    *gmap.AnyAnyMap
	aborted bool
}

var _ IEvent = new(Event)

// NewEvent 创建事件
func NewEvent(name string, data map[interface{}]interface{}) *Event {
	e := &Event{
		name: name,
	}
	if data == nil {
		e.data = gmap.NewAnyAnyMap(true)
	} else {
		e.data = gmap.NewAnyAnyMapFrom(data, true)
	}

	return e
}

// Get 获取元素
func (that *Event) Get(key interface{}) interface{} {
	if that.data == nil {
		return nil
	}
	return that.data.Get(key)
}

// Set 设置元素
func (that *Event) Set(key interface{}, val interface{}) {
	if that.data == nil {
		that.data = gmap.NewAnyAnyMap(true)
	}
	that.data.Set(key, val)
}

// Name 获取事件名
func (that *Event) Name() string {
	return that.name
}

// SetName 设置事件名
func (that *Event) SetName(name string) {
	that.name = name
}

// Data 获取事件的全部参数
func (that *Event) Data() map[interface{}]interface{} {
	if that.data == nil {
		return nil
	}
	return that.data.Map()
}

// SetData 设置事件的全部参数
func (that *Event) SetData(data map[interface{}]interface{}) IEvent {
	if that.data == nil {
		that.data = gmap.NewAnyAnyMap(true)
	}
	if data != nil {
		that.data.Sets(data)
	}
	return that
}

// Abort 设置事件是否终止
func (that *Event) Abort(abort bool) {
	that.aborted = abort
}

// IsAborted 判断事件是否终止
func (that *Event) IsAborted() bool {
	return that.aborted
}

// AttachTo 把事件加入到指定的管理器中
func (that *Event) AttachTo(m *EventBus) error {
	return m.AddEvent(that)
}
