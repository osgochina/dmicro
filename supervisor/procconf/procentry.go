package procconf

import (
	"github.com/gogf/gf/util/gconv"
	"os/exec"
	"strconv"
)

type ProcEntry struct {
	//进程名称
	Name string
	// 启动命令
	Command string
	// 启动参数
	Args []string

	//进程运行目录
	Directory string
	//在supervisord启动的时候也自动启动
	AutoStart bool
	//启动10秒后没有异常退出，就表示进程正常启动了，默认为1秒
	StartSecs int
	//程序退出后自动重启,可选值：[unexpected,true,false]，默认为unexpected，表示进程意外杀死后才重启
	autoReStart string
	// 进程退出的code值
	exitCodes string
	//启动失败自动重试次数，默认是3
	StartRetries int
	//进程重启间隔秒数，默认是0，表示不间隔
	RestartPause int
	//用哪个用户启动进程，默认是root
	User string
	//进程启动优先级，默认999，值小的优先启动
	Priority int

	//日志文件，需要注意当指定目录不存在时无法正常启动，所以需要手动创建目录（supervisord 会自动创建日志文件）
	stdoutLogfile string
	//stdout 日志文件大小，默认50MB
	stdoutLogFileMaxBytes string
	//stdout 日志文件备份数，默认是10
	stdoutLogFileBackups int
	// 把stderr重定向到stdout，默认false
	redirectStderr bool
	// 日志文件，进程启动后的标准错误写入该文件
	stderrLogfile string
	//stderr 日志文件大小，默认50MB
	stderrLogFileMaxBytes string
	//stderr 日志文件备份数，默认是10
	stderrLogFileBackups int
	//默认为false,进程被杀死时，是否向这个进程组发送stop信号，包括子进程
	StopAsGroup bool
	//默认为false，向进程组发送kill信号，包括子进程
	KillAsGroup bool
	//结束进程发送的信号
	StopSignal string
	// 发送结束进程的信号后等待的秒数
	StopWaitSecs int
	//
	KillWaitSecs int
	// 环境变量
	Environment map[string]string
	//当进程的二进制文件有修改，是否需要重启
	//RestartWhenBinaryChanged bool
	// 扩展配置
	extend map[string]string
}

func NewProcEntry() *ProcEntry {
	return &ProcEntry{
		AutoStart:             true,
		StartSecs:             1,
		autoReStart:           "true",
		StartRetries:          3,
		RestartPause:          0,
		User:                  "root",
		Priority:              999,
		stdoutLogfile:         "",
		stdoutLogFileMaxBytes: "50MB",
		stdoutLogFileBackups:  10,
		redirectStderr:        false,
		stderrLogfile:         "",
		stderrLogFileMaxBytes: "50MB",
		stderrLogFileBackups:  10,
		StopAsGroup:           false,
		KillAsGroup:           false,
		Environment:           make(map[string]string),
		//RestartWhenBinaryChanged:false,
		extend: make(map[string]string),
	}
}

func (that *ProcEntry) GetProcessName() string {
	return that.Name
}

func (that *ProcEntry) SetProcessName(name string) {
	that.Name = name
}

func (that *ProcEntry) CreateCommand() (*exec.Cmd, error) {
	return createCommand(that.Command, that.Args)
}

func (that *ProcEntry) SetStdoutLogFileMaxBytes(v string) {
	that.stdoutLogFileMaxBytes = v
}

func (that *ProcEntry) GetStdoutLogFileMaxBytes(defaultVal int) int {
	return that.getBytes(that.stdoutLogFileMaxBytes, defaultVal)
}

func (that *ProcEntry) SetStdoutLogFileBackups(v int) {
	that.stdoutLogFileBackups = v
}

func (that *ProcEntry) GetStdoutLogFileBackups(defaultVal int) int {
	if that.stdoutLogFileBackups > 0 {
		return that.stdoutLogFileBackups
	}
	return defaultVal
}

func (that *ProcEntry) SetStderrLogfile(v string) {
	that.stderrLogfile = v
}

func (that *ProcEntry) GetStderrLogfile(defaultVal string) string {
	if len(that.stderrLogfile) > 0 {
		return that.stderrLogfile
	}
	return defaultVal
}

func (that *ProcEntry) SetStdoutLogfile(v string) {
	that.stdoutLogfile = v
}

func (that *ProcEntry) GetStdoutLogfile(defaultVal string) string {
	if len(that.stdoutLogfile) > 0 {
		return that.stdoutLogfile
	}
	return defaultVal
}
func (that *ProcEntry) SetRedirectStderr(v bool) {
	that.redirectStderr = v
}

func (that *ProcEntry) GetRedirectStderr() bool {
	return that.redirectStderr
}

func (that *ProcEntry) SetStderrLogFileMaxBytes(v string) {
	that.stderrLogFileMaxBytes = v
}

func (that *ProcEntry) GetStderrLogFileMaxBytes(defaultVal int) int {
	return that.getBytes(that.stderrLogFileMaxBytes, defaultVal)
}

func (that *ProcEntry) SetStderrLogFileBackups(v int) {
	that.stderrLogFileBackups = v
}

func (that *ProcEntry) GetStderrLogFileBackups(defaultVal int) int {
	if that.stderrLogFileBackups > 0 {
		return that.stderrLogFileBackups
	}
	return defaultVal
}

func (that *ProcEntry) GetStopWaitSecs(defaultVal int) int {
	if that.StopWaitSecs > 0 {
		return that.StopWaitSecs
	}
	return defaultVal
}
func (that *ProcEntry) GetKillWaitSecs(defaultVal int) int {
	if that.KillWaitSecs > 0 {
		return that.KillWaitSecs
	}
	return defaultVal
}

// GetAutoReStart 是否自动重启
func (that *ProcEntry) GetAutoReStart(defaultVal string) string {
	if len(that.autoReStart) > 0 {
		return that.autoReStart
	}
	return defaultVal
}

// SetAutoReStart 设置自动重启 值为 true，false，unexpected
func (that *ProcEntry) SetAutoReStart(val string) {
	that.autoReStart = val
}

func (that *ProcEntry) GetExitCodes(defaultVal string) string {
	if len(that.exitCodes) > 0 {
		return that.exitCodes
	}
	return defaultVal
}

func (that *ProcEntry) SetExitCodes(val string) {
	that.exitCodes = val
}

// GetBytes returns value of the key as bytes setting.
//
//	logSize=1MB
//	logSize=1GB
//	logSize=1KB
//	logSize=1024
//
func (that *ProcEntry) getBytes(value string, defValue int) int {

	if len(value) > 2 {
		lastTwoBytes := value[len(value)-2:]
		if lastTwoBytes == "MB" {
			return that.toInt(value[:len(value)-2], 1024*1024, defValue)
		} else if lastTwoBytes == "GB" {
			return that.toInt(value[:len(value)-2], 1024*1024*1024, defValue)
		} else if lastTwoBytes == "KB" {
			return that.toInt(value[:len(value)-2], 1024, defValue)
		}
		return that.toInt(value, 1, defValue)
	}
	return defValue
}

func (that *ProcEntry) toInt(s string, factor int, defValue int) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i * factor
	}
	return defValue
}

func (that *ProcEntry) GetExtendInt(key string, defValue int) int {
	s, ok := that.extend[key]
	if ok {
		return gconv.Int(s)
	}
	return defValue
}

func (that *ProcEntry) GetExtendString(key string, defValue string) string {
	s, ok := that.extend[key]
	if ok {
		return s
	}
	return defValue
}
