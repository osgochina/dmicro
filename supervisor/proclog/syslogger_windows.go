package proclog

type SysLogger struct {
	NullLogger
}

// NewSysLogger 获取系统syslog的对象
func NewSysLogger(_ string, _ map[string]string) *SysLogger {

	logger := &SysLogger{}
	return logger
}

// NewRemoteSysLogger 获取远程系统日志的对象
func NewRemoteSysLogger(_ string, _ string, _ map[string]string) *SysLogger {
	logger := &SysLogger{}
	return logger

}

func (that *SysLogger) Write(_ []byte) (int, error) {
	return 0, nil
}

func (that *SysLogger) Close() error {
	return nil
}
