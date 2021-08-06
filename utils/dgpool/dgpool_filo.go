package dgpool

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"
)

const (
	// DefaultMaxGoroutinesAmount 默认的协程最大数量
	DefaultMaxGoroutinesAmount = 256 * 1024
	// DefaultMaxGoroutineIdleDuration 默认协程最大空闲时间
	DefaultMaxGoroutineIdleDuration = 10 * time.Second
)

type FILOPool struct {
	maxGoroutinesAmount      int              //最大的协程数量
	maxGoroutineIdleDuration time.Duration    // 最大空闲的时间
	lock                     sync.Mutex       //锁
	goroutinesCount          int              // 当前在执行的协程数量
	mustStop                 bool             // 强制停止
	ready                    []*goroutineChan // 准备好的可执行的协程
	stopCh                   chan struct{}    //协程池停止通知
	goroutineChanPool        sync.Pool        // 对象池
}

type goroutineChan struct {
	lastUseTime time.Time
	ch          chan func()
}

// NewFILOPool 创建协程栈池
func NewFILOPool(maxGoroutinesAmount int, maxGoroutineIdleDuration time.Duration) *FILOPool {
	fp := new(FILOPool)
	if maxGoroutinesAmount <= 0 {
		fp.maxGoroutinesAmount = DefaultMaxGoroutinesAmount
	} else {
		fp.maxGoroutinesAmount = maxGoroutinesAmount
	}
	if maxGoroutineIdleDuration <= 0 {
		fp.maxGoroutineIdleDuration = DefaultMaxGoroutineIdleDuration
	} else {
		fp.maxGoroutineIdleDuration = maxGoroutineIdleDuration
	}
	fp.start()
	return fp
}

func (that *FILOPool) MaxGoroutinesAmount() int {
	return that.maxGoroutinesAmount
}

func (that *FILOPool) MaxGoroutineIdle() time.Duration {
	return that.maxGoroutineIdleDuration
}

// 启动该协程池
func (that *FILOPool) start() {
	if that.stopCh != nil {
		panic("BUG: FILOPool already started")
	}
	that.stopCh = make(chan struct{})
	stopCh := that.stopCh
	go func() {
		var scratch []*goroutineChan
		for {
			// 把池子中的方法按顺序执行
			that.clean(&scratch)
			select {
			//如果协程池子已关闭，则结束
			case <-stopCh:
				return
			default:
				time.Sleep(that.maxGoroutineIdleDuration) // 空闲状态则暂停改协程
			}
		}
	}()
}

func (that *FILOPool) Stop() {
	if that.stopCh == nil {
		panic("BUG: FILOPool wasn't started")
	}
	close(that.stopCh)
	that.stopCh = nil

	// Stop all the goroutines waiting for incoming materiels.
	// Handle not wait for busy goroutines - they will stop after
	// serving the materiel and noticing gp.mustStop = true.
	that.lock.Lock()
	ready := that.ready
	for i, ch := range ready {
		ch.ch <- nil
		ready[i] = nil
	}
	that.ready = ready[:0]
	that.mustStop = true
	that.lock.Unlock()
}

// 清除池子中过时的协程
func (that *FILOPool) clean(scratch *[]*goroutineChan) {
	maxGoroutineIdleDuration := that.maxGoroutineIdleDuration

	//
	currentTime := time.Now()

	that.lock.Lock()
	ready := that.ready
	//有多少个协程准备好了
	n := len(ready)
	i := 0
	//如果池子中的任务还没有执行完毕，则跳过
	for i < n && currentTime.Sub(ready[i].lastUseTime) > maxGoroutineIdleDuration {
		i++
	}
	//
	*scratch = append((*scratch)[:0], ready[:i]...)
	if i > 0 {
		m := copy(ready, ready[i:])
		for i = m; i < n; i++ {
			ready[i] = nil
		}
		that.ready = ready[:m]
	}
	that.lock.Unlock()

	//通知过时的协程停止
	tmp := *scratch
	for j, ch := range tmp {
		ch.ch <- nil
		tmp[j] = nil
	}
}

var ErrLack = errors.New("lack of goroutines, because exceeded maxGoroutinesAmount limit")

// Go 通过 goroutine 执行方法，如果返回值不为nil，则表示超过了最大执行个数
func (that *FILOPool) Go(fn func()) error {
	ch := that.getCh()
	if ch == nil {
		return ErrLack
	}
	ch.ch <- fn
	return nil
}

// TryGo 尝试通过goroutine执行方法，如果不成功，则同步执行
func (that *FILOPool) TryGo(fn func()) {
	if that.Go(fn) != nil {
		fn()
	}
}

// MustGo 强制执行方法，直到执行完毕，或者上下文取消
func (that *FILOPool) MustGo(fn func(), ctx ...context.Context) error {
	if len(ctx) == 0 {
		for that.Go(fn) != nil {
			runtime.Gosched()
		}
		return nil
	}
	c := ctx[0]
	for {
		select {
		case <-c.Done():
			return c.Err()
		default:
			if that.Go(fn) == nil {
				return nil
			}
			runtime.Gosched()
		}
	}
}

// 获取可用的协程通道
func (that *FILOPool) getCh() *goroutineChan {
	var ch *goroutineChan
	createGoroutine := false

	that.lock.Lock()
	ready := that.ready
	n := len(ready) - 1
	// 如果可用的协程通道不存在
	if n < 0 {
		if that.goroutinesCount < that.maxGoroutinesAmount {
			createGoroutine = true
			that.goroutinesCount++
		}
	} else {
		// 把最后一个协程取过来
		ch = ready[n]
		ready[n] = nil
		that.ready = ready[:n]
	}
	that.lock.Unlock()

	if ch == nil {
		//如果不是新建的协程通道
		if !createGoroutine {
			return nil
		}
		// 从协程池子中获取可用协程对象
		vch := that.goroutineChanPool.Get()
		if vch == nil {
			vch = &goroutineChan{
				ch: make(chan func(), goroutineChanCap),
			}
		}
		ch = vch.(*goroutineChan)
		go func() {
			that.goroutineFunc(ch)
			that.goroutineChanPool.Put(vch)
		}()
	}
	return ch
}

// 在协程中执行方法
func (that *FILOPool) goroutineFunc(ch *goroutineChan) {
	for fn := range ch.ch {
		if fn == nil {
			break
		}
		fn()
		if !that.release(ch) {
			break
		}
	}

	that.lock.Lock()
	that.goroutinesCount--
	that.lock.Unlock()
}

func (that *FILOPool) release(ch *goroutineChan) bool {
	ch.lastUseTime = time.Now()
	that.lock.Lock()
	if that.mustStop {
		that.lock.Unlock()
		return false
	}
	that.ready = append(that.ready, ch)
	that.lock.Unlock()
	return true
}

var goroutineChanCap = func() int {
	// Use blocking goroutineChan if GOMAXPROCS=1.
	// This immediately switches Go to GoroutineFunc, which results
	// in higher performance (under go1.5 at least).
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	// Use non-blocking goroutineChan if GOMAXPROCS>1,
	// since otherwise the Go caller (Acceptor) may lag accepting
	// new materiels if GoroutineFunc is CPU-bound.
	return 1
}()
