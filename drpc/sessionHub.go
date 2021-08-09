package drpc

import (
	"github.com/osgochina/dmicro/utils"
)

type SessionHub struct {
	// key: session id (ip, name and so on)
	// value: *session
	sessions utils.Map
}

//新建一个session容器对象
func newSessionHub() *SessionHub {
	chub := &SessionHub{
		sessions: utils.AtomicMap(),
	}
	return chub
}

//写入一个session对象到池子
func (that *SessionHub) set(sess *session) {
	_sess, loaded := that.sessions.LoadOrStore(sess.ID(), sess)
	if !loaded {
		return
	}
	that.sessions.Store(sess.ID(), sess)
	if oldSess := _sess.(*session); sess != oldSess {
		_ = oldSess.Close()
	}
}

//从池子中获取一个session对象
func (that *SessionHub) get(id string) (*session, bool) {
	_sess, ok := that.sessions.Load(id)
	if !ok {
		return nil, false
	}
	return _sess.(*session), true
}

//迭代session对象池
func (that *SessionHub) rangeCallback(fn func(*session) bool) {
	that.sessions.Range(func(key, value interface{}) bool {
		return fn(value.(*session))
	})
}

func (that *SessionHub) random() (*session, bool) {
	_, sess, exist := that.sessions.Random()
	if !exist {
		return nil, false
	}
	return sess.(*session), true
}

//池子的长度
func (that *SessionHub) len() int {
	return that.sessions.Len()
}

//删除指定的session对象
func (that *SessionHub) delete(id string) {
	that.sessions.Delete(id)
}
