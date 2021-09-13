package logger

import "sync"

type CompositeLogger struct {
	lock    sync.Mutex
	loggers []Logger
}

func NewCompositeLogger(loggers []Logger) *CompositeLogger {
	return &CompositeLogger{loggers: loggers}
}

func (that *CompositeLogger) AddLogger(logger Logger) {
	that.lock.Lock()
	defer that.lock.Unlock()
	that.loggers = append(that.loggers, logger)
}

func (that *CompositeLogger) RemoveLogger(logger Logger) {
	that.lock.Lock()
	defer that.lock.Unlock()
	for i, t := range that.loggers {
		if t == logger {
			that.loggers = append(that.loggers[:i], that.loggers[i+1:]...)
			break
		}
	}
}

func (that *CompositeLogger) Write(p []byte) (n int, err error) {
	that.lock.Lock()
	defer that.lock.Unlock()

	for i, logger := range that.loggers {
		if i == 0 {
			n, err = logger.Write(p)
		} else {
			_, _ = logger.Write(p)
		}
	}
	return
}

func (that *CompositeLogger) Close() (err error) {
	that.lock.Lock()
	defer that.lock.Unlock()

	for i, logger := range that.loggers {
		if i == 0 {
			err = logger.Close()
		} else {
			_ = logger.Close()
		}
	}
	return
}

func (that *CompositeLogger) SetPid(pid int) {
	that.lock.Lock()
	defer that.lock.Unlock()

	for _, logger := range that.loggers {
		logger.SetPid(pid)
	}
}

// ReadLog read log data from first logger in CompositeLogger pool
func (that *CompositeLogger) ReadLog(offset int64, length int64) (string, error) {
	return that.loggers[0].ReadLog(offset, length)
}

// ReadTailLog tail the log data from first logger in CompositeLogger pool
func (that *CompositeLogger) ReadTailLog(offset int64, length int64) (string, int64, bool, error) {
	return that.loggers[0].ReadTailLog(offset, length)
}

// ClearCurLogFile clear the first logger file in CompositeLogger pool
func (that *CompositeLogger) ClearCurLogFile() error {
	return that.loggers[0].ClearCurLogFile()
}

// ClearAllLogFile clear all the files of first logger in CompositeLogger pool
func (that *CompositeLogger) ClearAllLogFile() error {
	return that.loggers[0].ClearAllLogFile()
}
