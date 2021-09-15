package proclog

import (
	"io"
	"strings"
	"sync"
)

// Logger 日志的接口
type Logger interface {
	io.WriteCloser
	SetPid(pid int)
	ReadLog(offset int64, length int64) (string, error)
	ReadTailLog(offset int64, length int64) (string, int64, bool, error)
	ClearCurLogFile() error
	ClearAllLogFile() error
}

// NullLocker 假锁
type NullLocker struct{}

func NewNullLocker() *NullLocker {
	return &NullLocker{}
}
func (that *NullLocker) Lock()   {}
func (that *NullLocker) Unlock() {}

// NewLogger 新建日志对象
func NewLogger(programName string, logFile string, locker sync.Locker, maxBytes int64, backups int, props map[string]string) Logger {
	files := splitLogFile(logFile)
	loggers := make([]Logger, 0)
	for i, f := range files {
		var lr Logger
		if i == 0 {
			lr = createLogger(programName, f, locker, maxBytes, backups, props)
		} else {
			lr = createLogger(programName, f, NewNullLocker(), maxBytes, backups, props)
		}
		loggers = append(loggers, lr)
	}
	return NewCompositeLogger(loggers)
}

func splitLogFile(logFile string) []string {
	files := strings.Split(logFile, ",")
	for i, f := range files {
		files[i] = strings.TrimSpace(f)
	}
	return files
}

// 创建日志对象
func createLogger(programName string, logFile string, locker sync.Locker, maxBytes int64, backups int, props map[string]string) Logger {

	if logFile == "/dev/stdout" {
		return NewStdoutLogger()
	}
	if logFile == "/dev/stderr" {
		return NewStderrLogger()
	}
	if logFile == "/dev/null" {
		return NewNullLogger()
	}

	if logFile == "syslog" {
		return NewSysLogger(programName, props)
	}
	if strings.HasPrefix(logFile, "syslog") {
		fields := strings.Split(logFile, "@")
		fields[0] = strings.TrimSpace(fields[0])
		fields[1] = strings.TrimSpace(fields[1])
		if len(fields) == 2 && fields[0] == "syslog" {
			return NewRemoteSysLogger(programName, fields[1], props)
		}
	}
	if len(logFile) > 0 {
		return NewFileLogger(logFile, maxBytes, backups, locker)
	}
	return NewNullLogger()
}
