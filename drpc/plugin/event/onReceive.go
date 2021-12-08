package event

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/eventbus"
)

type OnReceive struct {
	ReadCtx drpc.ReadCtx
	aborted bool
}

var _ eventbus.IEvent = new(OnReceive)

func newOnReceive(readCtx drpc.ReadCtx) *OnReceive {
	return &OnReceive{
		ReadCtx: readCtx,
	}
}

func (that *OnReceive) Name() string {
	return OnReceiveEvent
}

func (that *OnReceive) Get(_ interface{}) interface{} {
	return nil
}

// Set 设置元素
func (that *OnReceive) Set(_ interface{}, _ interface{}) {

}

// Data 获取事件的全部参数
func (that *OnReceive) Data() map[interface{}]interface{} {
	return nil
}

// SetData 设置事件的全部参数
func (that *OnReceive) SetData(_ map[interface{}]interface{}) eventbus.IEvent {

	return that
}

// Abort 设置事件是否终止
func (that *OnReceive) Abort(abort bool) {
	that.aborted = abort
}

// IsAborted 判断事件是否终止
func (that *OnReceive) IsAborted() bool {
	return that.aborted
}
