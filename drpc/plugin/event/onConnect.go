package event

import (
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/eventbus"
)

type OnConnect struct {
	Session   drpc.EarlySession
	IsSuccess bool
	Err       error
	IsRedial  bool
	aborted   bool
}

var _ eventbus.IEvent = new(OnConnect)

func newOnConnect(isSuccess bool, sess drpc.EarlySession, isRedial bool, err error) *OnConnect {
	return &OnConnect{
		Session:   sess,
		IsSuccess: isSuccess,
		IsRedial:  isRedial,
		Err:       err,
	}
}

func (that *OnConnect) Name() string {
	return OnConnectEvent
}

func (that *OnConnect) Get(_ interface{}) interface{} {
	return nil
}

// Set 设置元素
func (that *OnConnect) Set(_ interface{}, _ interface{}) {

}

// Data 获取事件的全部参数
func (that *OnConnect) Data() map[interface{}]interface{} {
	return nil
}

// SetData 设置事件的全部参数
func (that *OnConnect) SetData(_ map[interface{}]interface{}) eventbus.IEvent {
	return that
}

// Abort 设置事件是否终止
func (that *OnConnect) Abort(abort bool) {
	that.aborted = abort
}

// IsAborted 判断事件是否终止
func (that *OnConnect) IsAborted() bool {
	return that.aborted
}
