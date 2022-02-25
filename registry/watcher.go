package registry

import "time"

// Watcher 服务监视器
type Watcher interface {
	// Next 是一个阻塞调用
	Next() (*Result, error)
	//Stop()
}

// Result 监视器返回对象
type Result struct {
	Action  string
	Service *Service
}

// EventType 服务事件类型
type EventType int

const (
	// Create 有一个新的服务登记
	Create EventType = iota
	// Delete 有一个服务取消登记
	Delete
	// Update 服务的信息更新了
	Update
)

func (t EventType) String() string {
	switch t {
	case Create:
		return "create"
	case Delete:
		return "delete"
	case Update:
		return "update"
	default:
		return "unknown"
	}
}

// Event 事件
type Event struct {
	// Id 注册的id
	Id string
	// Type 事件类型
	Type EventType
	// Timestamp 事件发生事件
	Timestamp time.Time
	// Service 注册的服务
	Service *Service
}
