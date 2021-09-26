package eventbus

// Wildcard 事件名称，通配符，表示所有
const Wildcard = "*"

const (
	// Min 最低优先级
	Min = -300
	// Low 低优先级
	Low = -200
	// BelowNormal 相对低优先级
	BelowNormal = -100
	// Normal 正常优先级
	Normal = 0
	// AboveNormal 稍高优先级
	AboveNormal = 100
	// High 高优先级
	High = 200
	// Max 最高优先级
	Max = 300
)

// 默认事件总线
var defaultEventBus = New("default")

// Subscribe 订阅事件
// listener: 事件监听器,支持通过Listen()方法获取监听事件的名称
// priority：监听优先级
func Subscribe(listener IListener, priority ...int) (err error) {
	return defaultEventBus.Subscribe(listener, priority...)
}

// UnSubscribe 取消事件订阅
// listener: 事件监听器,支持通过Listen()方法获取监听事件的名称
func UnSubscribe(listener IListener) {
	defaultEventBus.UnSubscribe(listener)
}

// Listen 监听事件
// name: 事件名
// listener：事件监听器
// priority：监听优先级
func Listen(name string, listener BaseListener, priority ...int) (err error) {
	return defaultEventBus.Listen(name, listener, priority...)
}

// UnListen 取消事件的监听
// name: 事件名
// listener：事件监听器
func UnListen(name string, listener BaseListener) {
	defaultEventBus.UnListen(name, listener)
}

// RemoveListenersByName 取消指定事件名的所有事件监听
// name: 事件名
func RemoveListenersByName(name string) {
	defaultEventBus.RemoveListenersByName(name)
}

// RemoveListeners 移除监听方法
func RemoveListeners(listener IListener) {
	defaultEventBus.RemoveListeners(listener)
}

// HasListeners 判断是否存在指定名称的监听方法
// name: 事件名
func HasListeners(name string) bool {
	return defaultEventBus.HasListeners(name)
}

// Listeners 获取所有监听列表
func Listeners() map[string]*ListenerQueue {
	return defaultEventBus.Listeners()
}

// ListenersByName 获取指定事件名下的监听列表
func ListenersByName(name string) *ListenerQueue {
	return defaultEventBus.ListenersByName(name)
}

// ListenersCount 获取指定事件名下的监听方法数量
func ListenersCount(name string) int {
	return defaultEventBus.ListenersCount(name)
}

// ListenedNames 获取监听的事件名列表
func ListenedNames() []string {
	return defaultEventBus.ListenedNames()
}

// AddEvent 添加事件到管理器
// e: 事件对象
func AddEvent(e IEvent) error {
	return defaultEventBus.AddEvent(e)
}

// GetEvent 获取指定事件名的事件对象
func GetEvent(name string) (IEvent, bool) {
	return defaultEventBus.GetEvent(name)
}

// HasEvent 判断事件是否存在
func HasEvent(name string) bool {
	return defaultEventBus.HasEvent(name)
}

// RemoveEvent 移除自定义事件
func RemoveEvent(name string) {
	defaultEventBus.RemoveEvent(name)
}

// RemoveEvents 移除所有自定义事件
func RemoveEvents() {
	defaultEventBus.RemoveEvents()
}

// Fire 使用事件名触发事件
// name: 事件名
// params:需要传递给事件的参数
func Fire(name string, params map[interface{}]interface{}) (e IEvent, err error) {
	return defaultEventBus.Fire(name, params)
}

// Publish 使用事件对象触发事件
// e：事件对象
func Publish(e IEvent) error {
	return defaultEventBus.Publish(e)
}

// PublishBatch 批量触发事件
// Usage:
// 	PublishBatch("name1", "name2", &MyEvent{})
func PublishBatch(es ...interface{}) (ers []error) {
	return defaultEventBus.PublishBatch(es...)
}

// AsyncPublish 异步触发事件
// e：事件对象
func AsyncPublish(e IEvent) {
	defaultEventBus.AsyncPublish(e)
}
