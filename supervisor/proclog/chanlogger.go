package proclog

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
)

type ChanLogger struct {
	channel chan []byte
}

func NewChanLogger(channel chan []byte) *ChanLogger {
	return &ChanLogger{channel: channel}
}

func (that *ChanLogger) SetPid(_ int) {
	// NOTHING TO DO
	return
}

func (that *ChanLogger) Write(p []byte) (int, error) {
	that.channel <- p
	return len(p), nil
}

func (that *ChanLogger) Close() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	close(that.channel)
	return nil
}

func (that *ChanLogger) ReadLog(_ int64, _ int64) (string, error) {
	return "", gerror.New("NO_FILE")
}

func (that *ChanLogger) ReadTailLog(_ int64, _ int64) (string, int64, bool, error) {
	return "", 0, false, gerror.New("NO_FILE")
}

func (that *ChanLogger) ClearCurLogFile() error {
	return fmt.Errorf("NoLog")
}

func (that *ChanLogger) ClearAllLogFile() error {
	return gerror.New("NO_FILE")
}
