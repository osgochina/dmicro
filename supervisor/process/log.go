package process

import "github.com/osgochina/dmicro/supervisor/logger"

// 创建标准输出日志
func (that *Process) createStdoutLogger() logger.Logger {
	logFile := that.GetStdoutLogfile()
	maxBytes := int64(that.procEntry.GetStdoutLogFileMaxBytes(50 * 1024 * 1024))
	backups := that.procEntry.GetStdoutLogFileBackups(10)
	//logEventEmitter := p.createStdoutLogEventEmitter()
	props := make(map[string]string)
	syslogFacility := that.procEntry.GetExtendString("syslog_facility", "")
	syslogTag := that.procEntry.GetExtendString("syslog_tag", "")
	syslogPriority := that.procEntry.GetExtendString("syslog_stdout_priority", "")

	if len(syslogFacility) > 0 {
		props["syslog_facility"] = syslogFacility
	}
	if len(syslogTag) > 0 {
		props["syslog_tag"] = syslogTag
	}
	if len(syslogPriority) > 0 {
		props["syslog_priority"] = syslogPriority
	}
	return logger.NewLogger(that.GetName(), logFile, logger.NewNullLocker(), maxBytes, backups, props)
}

// 创建标准错误日志
func (that *Process) createStderrLogger() logger.Logger {
	logFile := that.GetStderrLogfile()
	maxBytes := int64(that.procEntry.GetStderrLogFileMaxBytes(50 * 1024 * 1024))
	backups := that.procEntry.GetStderrLogFileBackups(10)
	//logEventEmitter := p.createStderrLogEventEmitter()
	props := make(map[string]string)
	syslogFacility := that.procEntry.GetExtendString("syslog_facility", "")
	syslogTag := that.procEntry.GetExtendString("syslog_tag", "")
	syslogPriority := that.procEntry.GetExtendString("syslog_stdout_priority", "")

	if len(syslogFacility) > 0 {
		props["syslog_facility"] = syslogFacility
	}
	if len(syslogTag) > 0 {
		props["syslog_tag"] = syslogTag
	}
	if len(syslogPriority) > 0 {
		props["syslog_priority"] = syslogPriority
	}

	return logger.NewLogger(that.GetName(), logFile, logger.NewNullLocker(), maxBytes, backups, props)
}
