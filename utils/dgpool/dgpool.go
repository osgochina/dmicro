package dgpool

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/glog"
)

// Pool Goroutine Pool
type Pool struct {
	limit  int
	count  *gtype.Int
	list   *glist.List
	closed *gtype.Bool
}

//默认创建一个协程池
var pool = New()

// New 创建协程池，limit限制最多能同时运行多少个工作协程
func New(limit ...int) *Pool {
	p := &Pool{
		limit:  -1,
		count:  gtype.NewInt(),
		list:   glist.New(true),
		closed: gtype.NewBool(),
	}
	if len(limit) > 0 && limit[0] > 0 {
		p.limit = limit[0]
	}
	return p
}

// Go 执行协程
func Go(fn func()) bool {
	if err := pool.Add(fn); err != nil {
		glog.Errorf("%s", err.Error())
		return false
	}
	return true
}

// Add 往默认协程池中添加jobs
func Add(f func()) error {
	return pool.Add(f)
}

// AddWithRecover 在默认协程池中执行方法，并且执行完成后，如果出错，则调用recover方法
func AddWithRecover(userFunc func(), recoverFunc ...func(err error)) error {
	return pool.AddWithRecover(userFunc, recoverFunc...)
}

// AddWithSyncFunc 在默认协程池中执行方法，执行完成后回调
func AddWithSyncFunc(useFunc func(), syncFunc func(bool)) error {
	return pool.AddWithSyncFunc(useFunc, syncFunc)
}

// Size 默认协程池的大小
func Size() int {
	return pool.Size()
}

// Jobs 默认协程池当前中有多少个任务需要执行
func Jobs() int {
	return pool.Jobs()
}

// Add 添加待执行的方法到协程池
func (that *Pool) Add(f func()) error {

	//如果协程池已关闭，则返回错误
	for that.closed.Val() {
		return errors.New("pool closed")
	}
	//并发安全的把方法指针传入双向链表
	that.list.PushFront(f)
	var n int
	for {
		//获取协程池中的协程协程总数
		n = that.count.Val()
		//如果协程池中的协程已经超出限制，则不执行任务，直接返回
		if that.limit != -1 && n >= that.limit {
			return nil
		}
		//协程池中的协程够用，则执行
		if that.count.Cas(n, n+1) {
			break
		}
	}
	//启动一个协程执行任务
	that.fork()
	return nil
}

// AddWithRecover 添加任务，并在任务执行出错的情况下，回调recoverFunc
func (that *Pool) AddWithRecover(useFunc func(), recoverFunc ...func(err error)) error {
	return that.Add(func() {
		defer func() {
			if err := recover(); err != nil {
				if len(recoverFunc) > 0 && recoverFunc[0] == nil {
					recoverFunc[0](errors.New(fmt.Sprintf("%v", err)))
				}
			}
		}()
		useFunc()
	})
}

// AddWithSyncFunc 执行成功后回调方法
func (that *Pool) AddWithSyncFunc(useFunc func(), syncFunc func(bool)) error {
	return that.Add(func() {
		defer func() {
			if err := recover(); err != nil {
				syncFunc(false)
			} else {
				syncFunc(true)
			}
		}()
		useFunc()
	})
}

//开始执行一个协程
func (that *Pool) fork() {
	go func() {
		defer that.count.Add(-1)

		var job interface{}

		for !that.closed.Val() {
			if job = that.list.PopBack(); job != nil {
				job.(func())()
			} else {
				return
			}
		}
	}()
}

// IsClosed 判断协程池是否关闭
func (that *Pool) IsClosed() bool {
	return that.closed.Val()
}

// Close 关闭协程池
func (that *Pool) Close() {
	that.closed.Set(true)
}

// Size 协程池中的协程数量
func (that *Pool) Size() int {
	return that.count.Val()
}

// Jobs 协程池中的待执行任务数量
func (that *Pool) Jobs() int {
	return that.list.Size()
}

// Cap 协程池最大能够启动多少个协程
func (that *Pool) Cap() int {
	return that.limit
}
