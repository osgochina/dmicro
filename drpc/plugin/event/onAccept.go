package event

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/eventbus"
)

type OnAccept struct {
	Session drpc.EarlySession
	aborted bool
}

var _ eventbus.IEvent = new(OnAccept)

func newOnAccept(sess drpc.EarlySession) *OnAccept {
	return &OnAccept{
		Session: sess,
	}
}

func (that *OnAccept) Name() string {
	return OnAcceptEvent
}

func (that *OnAccept) Get(key interface{}) interface{} {
	return nil
}

// Set 设置元素
func (that *OnAccept) Set(key interface{}, val interface{}) {
}

// Data 获取事件的全部参数
func (that *OnAccept) Data() map[interface{}]interface{} {
	return nil
}

// SetData 设置事件的全部参数
func (that *OnAccept) SetData(data map[interface{}]interface{}) eventbus.IEvent {
	return that
}

// Abort 设置事件是否终止
func (that *OnAccept) Abort(abort bool) {
	that.aborted = abort
}

// IsAborted 判断事件是否终止
func (that *OnAccept) IsAborted() bool {
	return that.aborted
}
