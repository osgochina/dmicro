package drpc

import "github.com/gogf/gf/container/gmap"

type SessionHub struct {
	// key: session id (ip, name and so on)
	// value: *session
	sessions *gmap.Map
}

//新建一个session容器对象
func newSessionHub() *SessionHub {
	chub := &SessionHub{
		sessions: gmap.New(true),
	}
	return chub
}

//写入一个session对象到池子
func (that *SessionHub) set(sess *session) {
	_sess, loaded := that.sessions.Search(sess.ID())
	that.sessions.Set(sess.ID(), sess)
	if !loaded {
		return
	}
	if oldSess := _sess.(*session); sess != oldSess {
		_ = oldSess.Close()
	}
}

//从池子中获取一个session对象
func (that *SessionHub) get(id string) (*session, bool) {
	_sess, ok := that.sessions.Search(id)
	if !ok {
		return nil, false
	}
	return _sess.(*session), true
}

//迭代session对象池
func (that *SessionHub) rangeCallback(fn func(*session) bool) {
	that.sessions.Iterator(func(key, value interface{}) bool {
		return fn(value.(*session))
	})
}

//池子的长度
func (that *SessionHub) len() int {
	return that.sessions.Size()
}

//删除指定的session对象
func (that *SessionHub) delete(id string) {
	that.sessions.Remove(id)
}
