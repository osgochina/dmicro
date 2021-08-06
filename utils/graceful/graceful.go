package graceful

import (
	"context"
	"github.com/gogf/gf/os/glog"
	"os"
	"sync"
	"time"
)

// MinShutdownTimeout 最小停止超时时间
const MinShutdownTimeout = 15 * time.Second

type Graceful struct {
	shutdownTimeout              time.Duration
	firstSweep                   func() error
	beforeExiting                func() error
	locker                       sync.Mutex
	signal                       chan os.Signal
	logger                       *glog.Logger
	inheritedEnv                 map[string]string
	inheritedProcFiles           []*os.File
	defaultInheritedProcFilesLen int
}

func NewGraceful() *Graceful {
	graceful := &Graceful{
		signal: make(chan os.Signal),
		firstSweep: func() error {
			return nil
		},
		beforeExiting: func() error {
			return nil
		},
		logger:             glog.New(),
		inheritedEnv:       make(map[string]string),
		inheritedProcFiles: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	graceful.defaultInheritedProcFilesLen = len(graceful.inheritedProcFiles)
	return graceful
}

// Shutdown 执行进程关闭任务
func (that *Graceful) Shutdown(timeout ...time.Duration) {
	defer os.Exit(0)
	that.logger.Infof("shutting down process...")
	that.contextExec(timeout, "shutdown", func(ctxTimeout context.Context) <-chan struct{} {
		endCh := make(chan struct{})
		go func() {
			defer close(endCh)

			var graceful = true

			if err := that.firstSweep(); err != nil {
				that.logger.Errorf("[shutdown-firstSweep] %s", err.Error())
				graceful = false
			}

			graceful = that.shutdown(ctxTimeout, "shutdown") && graceful

			if graceful {
				that.logger.Infof("process is shutdown gracefully!")
			} else {
				that.logger.Infof("process is shutdown, but not gracefully!")
			}
		}()
		return endCh
	})
}

// SetShutdown 设置退出的基本参数
func (that *Graceful) SetShutdown(timeout time.Duration, firstSweepFunc, beforeExitingFunc func() error) {
	if timeout < 0 {
		that.shutdownTimeout = 1<<63 - 1
	} else if timeout < MinShutdownTimeout {
		that.shutdownTimeout = MinShutdownTimeout
	} else {
		that.shutdownTimeout = timeout
	}
	if firstSweepFunc == nil {
		firstSweepFunc = func() error { return nil }
	}
	if beforeExitingFunc == nil {
		beforeExitingFunc = func() error { return nil }
	}
	that.firstSweep = firstSweepFunc
	that.beforeExiting = beforeExitingFunc
}

// 执行shutdown和reboot命令，并且计时，在规定的时候内为执行完收尾动作，则强制结束进程
func (that *Graceful) contextExec(timeout []time.Duration, action string, deferCallback func(ctxTimeout context.Context) <-chan struct{}) {
	if len(timeout) > 0 {
		that.SetShutdown(timeout[0], that.firstSweep, that.beforeExiting)
	}

	ctxTimeout, _ := context.WithTimeout(context.Background(), that.shutdownTimeout)
	select {
	case <-ctxTimeout.Done():
		if err := ctxTimeout.Err(); err != nil {
			that.logger.Errorf("[%s-timeout] %s", action, err.Error())
		}
	case <-deferCallback(ctxTimeout):
	}
}

//执行后置函数
func (that *Graceful) shutdown(ctxTimeout context.Context, action string) bool {
	if err := that.beforeExiting(); err != nil {
		that.logger.Errorf("[%s-beforeExiting] %s", action, err.Error())
		return false
	}
	return true
}
