package eventbus

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
)

// BaseListener 基础监听器
type BaseListener interface {
	Process(event IEvent) error
}

// IListener 完全体监听器
type IListener interface {
	Listen() []interface{}
	Process(event IEvent) error
}

// ListenerFunc 强制把监听的方法转换成对象
type ListenerFunc func(e IEvent) error

func (fn ListenerFunc) Process(e IEvent) error {
	return fn(e)
}

// ListenerItem 事件监听的数据结构
type ListenerItem struct {
	Priority int
	Listener BaseListener
}

// NewListenerItem 新建事件监听条目
func NewListenerItem(listener BaseListener, priority int) *ListenerItem {
	return &ListenerItem{
		Priority: priority,
		Listener: listener,
	}
}

// ListenerQueue 监听队列，以优先级排序
type ListenerQueue struct {
	data *garray.SortedArray
}

// 创建监听队列
func newListenerQueue() *ListenerQueue {
	queue := &ListenerQueue{
		data: garray.NewSortedArray(comparator, true),
	}
	return queue
}

func comparator(a, b interface{}) int {
	if a.(*ListenerItem).Priority < b.(*ListenerItem).Priority {
		return 1
	}
	if a.(*ListenerItem).Priority > b.(*ListenerItem).Priority {
		return -1
	}
	return 0
}

// Len 获取优先队列的长度
func (that *ListenerQueue) Len() int {
	return that.data.Len()
}

// IsEmpty 判断监听队列是否为空
func (that *ListenerQueue) IsEmpty() bool {
	return that.data.IsEmpty()
}

// Add 存入条目到优先队列
func (that *ListenerQueue) Add(item *ListenerItem) {
	that.data.Add(item)
}

// Items 获取所有监听器条目
func (that *ListenerQueue) Items() []*ListenerItem {
	var items = make([]*ListenerItem, that.data.Len())
	for k, v := range that.data.Slice() {
		items[k] = v.(*ListenerItem)
	}
	return items
}

// Remove 移除指定的监听
func (that *ListenerQueue) Remove(listener BaseListener) {
	ptrVal := fmt.Sprintf("%p", listener)
	var newData *garray.SortedArray
	that.data.Iterator(func(k int, v interface{}) bool {
		liPtrVal := fmt.Sprintf("%p", v.(*ListenerItem).Listener)
		if liPtrVal == ptrVal {
			return true
		}
		if newData == nil {
			newData = garray.NewSortedArray(comparator, true)
		}
		newData.Add(v)
		return true
	})
	if newData != nil {
		that.data = newData
	}
}

// Clear 清除事件
func (that *ListenerQueue) Clear() {
	that.data.Clear()
}
