package event

import (
	"fmt"
)

// ListenerItem 事件监听的数据结构
type ListenerItem struct {
	EventName string
	Priority  int
	Listener  IListener
}

// NewListenerItem 新建事件监听条目
func NewListenerItem(name string, listener IListener, priority int) *ListenerItem {
	return &ListenerItem{
		EventName: name,
		Priority:  priority,
		Listener:  listener,
	}
}

// ListenerQueue 事件监听队列
type ListenerQueue struct {
	items []*ListenerItem
}

// NewListenerQueue 创建监听队列
func newListenerQueue() *ListenerQueue {
	return &ListenerQueue{}
}

// Len 获取监听队列的长度
func (that *ListenerQueue) Len() int {
	return len(that.items)
}

// IsEmpty 判断监听队列是否为空
func (that *ListenerQueue) IsEmpty() bool {
	return len(that.items) == 0
}

// Push 把监听事件放入队列
func (that *ListenerQueue) Push(li *ListenerItem) *ListenerQueue {
	that.items = append(that.items, li)
	return that
}

// GetListenersForEvent 传入事件名，获取优先队列
func (that *ListenerQueue) GetListenersForEvent(name string) *PriorityQueue {
	var queue *PriorityQueue
	for _, item := range that.items {
		if item.EventName == name {
			if queue == nil {
				queue = newPriorityQueue()
			}
			queue.Push(item)
		}
	}
	return queue
}

// Remove 移除指定的监听
func (that *ListenerQueue) Remove(listener IListener) {
	if listener == nil {
		return
	}

	// unsafe.Pointer(listener)
	ptrVal := fmt.Sprintf("%p", listener)

	var newItems []*ListenerItem
	for _, li := range that.items {
		liPtrVal := fmt.Sprintf("%p", li.Listener)
		if liPtrVal == ptrVal {
			continue
		}

		newItems = append(newItems, li)
	}

	that.items = newItems
}

// Clear 清除事件
func (that *ListenerQueue) Clear() {
	that.items = that.items[:0]
}
