package eventbus

import (
	"container/heap"
	"fmt"
	"sync"
)

// BaseListener 基础监听器
type BaseListener interface {
	Process(event IEvent) error
}

// Listener 完全体监听器
type Listener interface {
	Listen() []string
	Process(event IEvent) error
}

// ListenerFunc 强制把监听的方法转换成对象
type ListenerFunc func(e IEvent) error

func (fn ListenerFunc) Process(e IEvent) error {
	return fn(e)
}

// ListenerItem 事件监听的数据结构
type ListenerItem struct {
	EventName string
	Priority  int
	Listener  BaseListener
}

// NewListenerItem 新建事件监听条目
func NewListenerItem(name string, listener BaseListener, priority int) *ListenerItem {
	return &ListenerItem{
		EventName: name,
		Priority:  priority,
		Listener:  listener,
	}
}

// ListenerQueue 监听队列，以优先级排序
type ListenerQueue struct {
	mu   sync.RWMutex
	heap *priorityQueueHeap
}

// 创建监听队列
func newListenerQueue() *ListenerQueue {
	queue := &ListenerQueue{
		heap: &priorityQueueHeap{
			array: make([]*ListenerItem, 0),
		},
	}
	heap.Init(queue.heap)
	return queue
}

// Len 获取优先队列的长度
func (that *ListenerQueue) Len() int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return that.heap.Len()
}

// IsEmpty 判断监听队列是否为空
func (that *ListenerQueue) IsEmpty() bool {
	return that.Len() == 0
}

// Push 存入条目到优先队列
func (that *ListenerQueue) Push(item *ListenerItem) {
	that.mu.Lock()
	defer that.mu.Unlock()
	heap.Push(that.heap, item)
}

// Pop 从优先队列中弹出条目
func (that *ListenerQueue) Pop() *ListenerItem {
	that.mu.Lock()
	defer that.mu.Unlock()
	if v := heap.Pop(that.heap); v != nil {
		item := v.(*ListenerItem)
		return item
	}
	return nil
}

// Remove 移除指定的监听
func (that *ListenerQueue) Remove(listener BaseListener) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if listener == nil {
		return
	}
	// unsafe.Pointer(listener)
	ptrVal := fmt.Sprintf("%p", listener)

	var newItems []*ListenerItem
	for _, li := range that.heap.array {
		liPtrVal := fmt.Sprintf("%p", li.Listener)
		if liPtrVal == ptrVal {
			continue
		}

		newItems = append(newItems, li)
	}

	that.heap.array = newItems
	heap.Init(that.heap)
}

// Clear 清除事件
func (that *ListenerQueue) Clear() {
	that.heap.array = that.heap.array[:0]
}

///--------------------------------------------优先队列的实现----------------------------------------------/////

// 优先队列的实现
type priorityQueueHeap struct {
	array []*ListenerItem
}

// Len 用来实现接口 sort.Interface.
func (that *priorityQueueHeap) Len() int {
	return len(that.array)
}

// Less 用来实现接口 sort.Interface.
func (that *priorityQueueHeap) Less(i, j int) bool {
	return that.array[i].Priority < that.array[j].Priority
}

// Swap 用来实现接口 sort.Interface.
func (that *priorityQueueHeap) Swap(i, j int) {
	if len(that.array) == 0 {
		return
	}
	that.array[i], that.array[j] = that.array[j], that.array[i]
}

// Push 写入数据到堆.
func (that *priorityQueueHeap) Push(x interface{}) {
	that.array = append(that.array, x.(*ListenerItem))
}

// Pop 从堆末尾返回一条数据
func (that *priorityQueueHeap) Pop() interface{} {
	length := len(that.array)
	if length == 0 {
		return nil
	}
	item := that.array[length-1]
	that.array = that.array[0 : length-1]
	return item
}
