package event

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/eventbus"
)

type OnClose struct {
	Session drpc.BaseSession
	aborted bool
}

var _ eventbus.IEvent = new(OnClose)

func newOnClose(sess drpc.BaseSession) *OnClose {
	return &OnClose{
		Session: sess,
	}
}

func (that *OnClose) Name() string {
	return OnCloseEvent
}

func (that *OnClose) Get(_ interface{}) interface{} {
	return nil
}

// Set 设置元素
func (that *OnClose) Set(_ interface{}, _ interface{}) {
}

// Data 获取事件的全部参数
func (that *OnClose) Data() map[interface{}]interface{} {
	return nil
}

// SetData 设置事件的全部参数
func (that *OnClose) SetData(_ map[interface{}]interface{}) eventbus.IEvent {
	return that
}

// Abort 设置事件是否终止
func (that *OnClose) Abort(abort bool) {
	that.aborted = abort
}

// IsAborted 判断事件是否终止
func (that *OnClose) IsAborted() bool {
	return that.aborted
}
