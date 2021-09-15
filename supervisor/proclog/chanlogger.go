package proclog

import (
	"fmt"
	"github.com/gogf/gf/errors/gerror"
)

type ChanLogger struct {
	channel chan []byte
}

func NewChanLogger(channel chan []byte) *ChanLogger {
	return &ChanLogger{channel: channel}
}

func (that *ChanLogger) SetPid(pid int) {
	// NOTHING TO DO
}

func (that *ChanLogger) Write(p []byte) (int, error) {
	that.channel <- p
	return len(p), nil
}

func (that *ChanLogger) Close() error {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	close(that.channel)
	return nil
}

func (that *ChanLogger) ReadLog(offset int64, length int64) (string, error) {
	return "", gerror.New("NO_FILE")
}

func (that *ChanLogger) ReadTailLog(offset int64, length int64) (string, int64, bool, error) {
	return "", 0, false, gerror.New("NO_FILE")
}

func (that *ChanLogger) ClearCurLogFile() error {
	return fmt.Errorf("NoLog")
}

func (that *ChanLogger) ClearAllLogFile() error {
	return gerror.New("NO_FILE")
}
