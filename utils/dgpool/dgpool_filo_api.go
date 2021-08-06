package dgpool

import (
	"context"
	"github.com/gogf/gf/os/glog"
	"time"
)

var (
	_maxGoroutinesAmount      = (1024 * 1024 * 8) / 8 // max memory 8GB (8KB/goroutine)
	_maxGoroutineIdleDuration time.Duration
	_filoPool                 = NewFILOPool(_maxGoroutinesAmount, _maxGoroutineIdleDuration)
)

func SetFILOPool(maxGoroutinesAmount int, maxGoroutineIdleDuration time.Duration) {
	_maxGoroutinesAmount, _maxGoroutineIdleDuration := maxGoroutinesAmount, maxGoroutineIdleDuration
	if _filoPool != nil {
		_filoPool.Stop()
	}
	_filoPool = NewFILOPool(_maxGoroutinesAmount, _maxGoroutineIdleDuration)
}

// FILOGo 使用栈的形式组织协程执行方法，
func FILOGo(fn func()) bool {
	if err := _filoPool.Go(fn); err != nil {
		glog.Printf("%s", err.Error())
		return false
	}
	return true
}

// FILOAnywayGo 强制执行方法
func FILOAnywayGo(fn func()) {
	_filoPool.MustGo(fn)
}

// FILOMustGo 强制执行方法，并且传入上下文
func FILOMustGo(fn func(), ctx ...context.Context) error {
	return _filoPool.MustGo(fn, ctx...)
}

// FILOTryGo 尝试执行方法
func FILOTryGo(fn func()) {
	_filoPool.TryGo(fn)
}
