package proclog

import (
	"io"
	"os"
)

type StdLogger struct {
	NullLogger
	//logEventEmitter LogEventEmitter
	writer io.Writer
}

func NewStdoutLogger() *StdLogger {
	return &StdLogger{
		writer: os.Stdout,
	}
}
func NewStderrLogger() *StdLogger {
	return &StdLogger{
		writer: os.Stderr,
	}
}

func (that *StdLogger) Write(p []byte) (int, error) {
	n, err := that.writer.Write(p)
	return n, err
}
