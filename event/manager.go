package event

import (
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gregex"
	"strings"
	"sync"
)

// Wildcard 事件名称，通配符，表示所有
const Wildcard = "*"

type Manager struct {
	sync.Mutex
	EnableLock    bool
	name          string
	events        map[string]IEvent
	listeners     *ListenerQueue
	listenedNames *gset.StrSet
}

func NewManager(name string) *Manager {
	em := &Manager{
		name:   name,
		events: make(map[string]IEvent),
		// listeners
		listeners:     newListenerQueue(),
		listenedNames: gset.NewStrSet(true),
	}
	return em
}

// AddListener 添加监听事件
func (that *Manager) AddListener(name string, listener IListener, priority ...int) error {
	return that.On(name, listener, priority...)
}

// Listen On的别名
func (that *Manager) Listen(name string, listener IListener, priority ...int) (err error) {
	return that.On(name, listener, priority...)
}

// On 监听事件
func (that *Manager) On(name string, listener IListener, priority ...int) (err error) {
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
	that.listenedNames.Add(name)
	that.listeners.Push(NewListenerItem(name, listener, pv))
	return nil
}

// AddEvent 添加自定义事件
func (that *Manager) AddEvent(e IEvent) error {
	name, err := checkName(e.Name())
	if err != nil {
		return err
	}
	that.events[name] = e
	return nil
}

// GetEvent 根据name获取自定义事件
func (that *Manager) GetEvent(name string) (e IEvent, ok bool) {
	e, ok = that.events[name]
	return
}

// HasEvent 判断自定义事件是否存在
func (that *Manager) HasEvent(name string) bool {
	_, ok := that.events[name]
	return ok
}

// RemoveEvent 移除自定义事件
func (that *Manager) RemoveEvent(name string) {
	if _, ok := that.events[name]; ok {
		delete(that.events, name)
	}
}

// RemoveEvents 移除所有自定义事件
func (that *Manager) RemoveEvents() {
	that.events = map[string]IEvent{}
}

// 检查事件名是否符合规范
func checkName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", gerror.New("event: 事件名不能为空")
	}

	if !gregex.IsMatchString(`^[a-zA-Z][\w-.*]*$`, name) {
		return "", gerror.New(`event: 事件名格式不正确,请匹配'^[a-zA-Z][\w-.]*$'`)
	}

	return name, nil
}

// HasListeners 判断是否存在指定名称的监听
func (that *Manager) HasListeners(name string) bool {
	return that.listenedNames.Contains(name)
}

// Clear 清空事件管理对象
func (that *Manager) Clear() {
	that.Reset()
}

// Reset 重置事件管理对象
func (that *Manager) Reset() {
	// clear all listeners
	that.listeners.Clear()
	// reset all
	that.name = ""
	that.events = make(map[string]IEvent)
	that.listeners = newListenerQueue()
	that.listenedNames.Clear()
}
