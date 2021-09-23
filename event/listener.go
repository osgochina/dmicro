package event

type IListener interface {
	Listen() []IEvent
	Process(event IEvent) error
}

// ListenerFunc 强制把监听的方法转换成对象
type ListenerFunc func(e IEvent) error

func (fn ListenerFunc) Listen() []IEvent {
	return []IEvent{}
}

func (fn ListenerFunc) Process(e IEvent) error {
	return fn(e)
}
