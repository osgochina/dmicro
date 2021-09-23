package event

import (
	"container/heap"
	"sync"
)

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

type PriorityQueue struct {
	mu   sync.RWMutex
	heap *priorityQueueHeap
}

type priorityQueueHeap struct {
	array []*ListenerItem
}

func newPriorityQueue() *PriorityQueue {
	queue := &PriorityQueue{
		heap: &priorityQueueHeap{
			array: make([]*ListenerItem, 0),
		},
	}
	heap.Init(queue.heap)
	return queue
}

// Len 获取优先队列的长度
func (that *PriorityQueue) Len() int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return that.heap.Len()
}

// Push 存入条目到优先队列
func (that *PriorityQueue) Push(item *ListenerItem) {
	that.mu.Lock()
	defer that.mu.Unlock()
	heap.Push(that.heap, item)
}

// Pop 从优先队列中弹出条目
func (that *PriorityQueue) Pop() *ListenerItem {
	that.mu.Lock()
	defer that.mu.Unlock()
	if v := heap.Pop(that.heap); v != nil {
		item := v.(*ListenerItem)
		return item
	}
	return nil
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
