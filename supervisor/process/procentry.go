package process

import (
	"os/exec"
	"strconv"
	"syscall"
)

type Entry struct {
	//进程名称
	name string
	// 启动命令
	command string
	// 启动参数
	args []string

	//进程运行目录
	directory string
	//在supervisord启动的时候也自动启动
	autoStart bool
	//启动10秒后没有异常退出，就表示进程正常启动了，默认为1秒
	startSecs int
	//程序退出后自动重启,可选值：[unexpected,true,false]，默认为unexpected，表示进程意外杀死后才重启
	autoReStart string
	// 进程退出的code值
	exitCodes string
	//启动失败自动重试次数，默认是3
	startRetries int
	//进程重启间隔秒数，默认是0，表示不间隔
	restartPause int
	//用哪个用户启动进程，默认是root
	user string
	//进程启动优先级，默认999，值小的优先启动
	priority int

	//日志文件，需要注意当指定目录不存在时无法正常启动，所以需要手动创建目录（supervisord 会自动创建日志文件）
	stdoutLogfile string
	//stdout 日志文件大小，默认50MB
	stdoutLogFileMaxBytes int
	//stdout 日志文件备份数，默认是10
	stdoutLogFileBackups int
	// 把stderr重定向到stdout，默认false
	redirectStderr bool
	// 日志文件，进程启动后的标准错误写入该文件
	stderrLogfile string
	//stderr 日志文件大小，默认50MB
	stderrLogFileMaxBytes int
	//stderr 日志文件备份数，默认是10
	stderrLogFileBackups int
	//默认为false,进程被杀死时，是否向这个进程组发送stop信号，包括子进程
	stopAsGroup bool
	//默认为false，向进程组发送kill信号，包括子进程
	killAsGroup bool
	//结束进程发送的信号
	stopSignal string
	// 发送结束进程的信号后等待的秒数
	stopWaitSecs int
	// 强杀进程等待秒数
	killWaitSecs int
	// 环境变量
	environment []string
	//当进程的二进制文件有修改，是否需要重启
	//RestartWhenBinaryChanged bool

	// 扩展参数
	extend map[string]string
}

// NewEntry 创建进程启动配置
func NewEntry(command string, args ...[]string) *Entry {
	proc := &Entry{
		command:               command,
		name:                  command,
		autoStart:             true,
		startSecs:             1,
		autoReStart:           "true",
		startRetries:          3,
		restartPause:          0,
		user:                  "root",
		priority:              999,
		stdoutLogfile:         "",
		stdoutLogFileMaxBytes: 50 * 1024 * 1024,
		stdoutLogFileBackups:  10,
		redirectStderr:        false,
		stderrLogfile:         "",
		stderrLogFileMaxBytes: 50 * 1024 * 1024,
		stderrLogFileBackups:  10,
		stopAsGroup:           false,
		killAsGroup:           false,
		extend:                make(map[string]string),
	}
	if len(args) > 0 {
		proc.args = args[0]
	}
	return proc
}

// CreateCommand 根据就配置生成cmd对象
func (that *Entry) CreateCommand() (*exec.Cmd, error) {

	cmd := exec.Command(that.command)
	if len(that.args) > 0 {
		cmd.Args = append([]string{that.command}, that.args...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	return cmd, nil
}

// Name 获取进程名称
func (that *Entry) Name() string {
	return that.name
}

// SetName 设置进程名
func (that *Entry) SetName(name string) {
	that.name = name
}

// Command 获取启动命令
func (that *Entry) Command() string {
	return that.command
}

// SetCommand 设置启动命令
func (that *Entry) SetCommand(command string) {
	that.command = command
}

// Args 获取参数
func (that *Entry) Args() []string {
	return that.args
}

// SetArgs 设置参数
func (that *Entry) SetArgs(args []string) {
	that.args = args
}

// Directory 程序运行目录
func (that *Entry) Directory() string {
	return that.directory
}

// SetDirectory 设置程序运行目录
func (that *Entry) SetDirectory(directory string) {
	that.directory = directory
}

// User 启动进程的用户
func (that *Entry) User() string {
	return that.user
}

// SetUser 设置启动进程的用户
func (that *Entry) SetUser(user string) {
	that.user = user
}

// AutoStart 判断程序是否需要自动启动
func (that *Entry) AutoStart() bool {
	return that.autoStart
}

// SetAutoStart 设置程序是否需要自动启动
func (that *Entry) SetAutoStart(autoStart bool) {
	that.autoStart = autoStart
}

// StartSecs 指定启动多少秒后没有异常退出，则表示启动成功
func (that *Entry) StartSecs() int {
	return that.startSecs
}

// SetStartSecs 指定启动多少秒后没有异常退出，则表示启动成功
// 未设置该值，则表示cmd.Start方法调用为出错，则表示启动成功，
// 设置了该值，则表示程序启动后需稳定运行指定的秒数后才算启动成功
func (that *Entry) SetStartSecs(startSecs int) {
	that.startSecs = startSecs
}

// AutoReStart 值为 true，false，unexpected
func (that *Entry) AutoReStart(defVal string) string {
	if len(that.autoReStart) > 0 {
		return that.autoReStart
	}
	return defVal
}

// SetAutoReStart 设置自动重启 值为 true，false，unexpected
func (that *Entry) SetAutoReStart(val string) {
	that.autoReStart = val
}

// ExitCodes 退出code的值列表，格式为 1，2，3，4
func (that *Entry) ExitCodes(defaultVal string) string {
	if len(that.exitCodes) > 0 {
		return that.exitCodes
	}
	return defaultVal
}

// SetExitCodes 设置退出code的值列表，格式为 1，2，3，4
func (that *Entry) SetExitCodes(val string) {
	that.exitCodes = val
}

// StartRetries 启动失败自动重试次数
func (that *Entry) StartRetries() int {
	return that.startRetries
}

// SetStartRetries 设置启动失败自动重试次数
func (that *Entry) SetStartRetries(startRetries int) {
	that.startRetries = startRetries
}

// RestartPause 进程重启间隔秒数，默认是0，表示不间隔
func (that *Entry) RestartPause() int {
	return that.restartPause
}

// SetRestartPause 设置进程重启间隔秒数
func (that *Entry) SetRestartPause(restartPause int) {
	that.restartPause = restartPause
}

// Priority 进程启动的优先级
func (that *Entry) Priority() int {
	return that.priority
}

// SetPriority 设置进程启动优先级
func (that *Entry) SetPriority(priority int) {
	that.priority = priority
}

// StdoutLogfile 标准输出的log文件地址
func (that *Entry) StdoutLogfile(defaultVal string) string {
	if len(that.stdoutLogfile) > 0 {
		return that.stdoutLogfile
	}
	return defaultVal
}

// SetStdoutLogfile 设置标准输出的log文件地址
func (that *Entry) SetStdoutLogfile(v string) {
	that.stdoutLogfile = v
}

// StdoutLogFileMaxBytes 标准输出的log文件最大容量，大于这个容量会分包
func (that *Entry) StdoutLogFileMaxBytes(defaultVal int) int {
	if that.stdoutLogFileMaxBytes <= 0 {
		return defaultVal
	}
	return that.stdoutLogFileMaxBytes
}

// SetStdoutLogFileMaxBytes 设置标准输出的log文件最大容量，格式：5KB，10MB，2GB，默认是50MB
func (that *Entry) SetStdoutLogFileMaxBytes(v string) {
	that.stdoutLogFileMaxBytes = that.getBytes(v, 50*1024*1024)
}

// StdoutLogFileBackups 标准输出的log文件如果达到了最大容量，会自动分包，这个值是设置它的最大分包份数
func (that *Entry) StdoutLogFileBackups(defaultVal int) int {
	if that.stdoutLogFileBackups > 0 {
		return that.stdoutLogFileBackups
	}
	return defaultVal
}

// SetStdoutLogFileBackups 标准输出的log文件如果达到了最大容量，会自动分包，这个值是设置它的最大分包份数
func (that *Entry) SetStdoutLogFileBackups(v int) {
	that.stdoutLogFileBackups = v
}

// RedirectStderr 是否重写标准错误输出到标准输出
func (that *Entry) RedirectStderr() bool {
	return that.redirectStderr
}

// SetRedirectStderr 是否重写标准错误输出到标准输出
func (that *Entry) SetRedirectStderr(v bool) {
	that.redirectStderr = v
}

// StderrLogfile 标准错误输出的log文件地址
func (that *Entry) StderrLogfile(defaultVal string) string {
	if len(that.stderrLogfile) > 0 {
		return that.stderrLogfile
	}
	return defaultVal
}

// SetStderrLogfile 设置标准错误输出的log文件地址
func (that *Entry) SetStderrLogfile(v string) {
	that.stderrLogfile = v
}

// StderrLogFileMaxBytes 标准错误输出的log文件最大容量，大于这个容量会分包
func (that *Entry) StderrLogFileMaxBytes(defaultVal int) int {
	if that.stderrLogFileMaxBytes <= 0 {
		return defaultVal
	}
	return that.stderrLogFileMaxBytes
}

// SetStderrLogFileMaxBytes 设置标准错误输出的log文件最大容量，格式：5KB，10MB，2GB，默认是50MB
func (that *Entry) SetStderrLogFileMaxBytes(v string) {
	that.stderrLogFileMaxBytes = that.getBytes(v, 50*1024*1024)
}

// StderrLogFileBackups 标准错误输出的log文件如果达到了最大容量，会自动分包，这个值是设置它的最大分包份数
func (that *Entry) StderrLogFileBackups(defaultVal int) int {
	if that.stderrLogFileBackups > 0 {
		return that.stderrLogFileBackups
	}
	return defaultVal
}

// SetStderrLogFileBackups 设置标准错误输出的log文件如果达到了最大容量
func (that *Entry) SetStderrLogFileBackups(v int) {
	that.stderrLogFileBackups = v
}

// StopAsGroup 停止进程的时候，是否向该进程组发送停止信号
func (that *Entry) StopAsGroup() bool {
	return that.stopAsGroup
}

// SetStopAsGroup 停止进程的时候，是否向该进程组发送停止信号
func (that *Entry) SetStopAsGroup(stopAsGroup bool) {
	that.stopAsGroup = stopAsGroup
}

// KillAsGroup 强制杀死进程的时候，是否向该进程的进程组发送kill信号
func (that *Entry) KillAsGroup() bool {
	return that.killAsGroup
}

// SetKillAsGroup 强制杀死进程的时候，是否向该进程的进程组发送kill信号
func (that *Entry) SetKillAsGroup(killAsGroup bool) {
	that.killAsGroup = killAsGroup
}

// StopSignal 正常结束进程是否需要发送的信号
func (that *Entry) StopSignal() string {
	return that.stopSignal
}

// SetStopSignal 正常结束进程是否需要发送的信号
func (that *Entry) SetStopSignal(stopSignal string) {
	that.stopSignal = stopSignal
}

// StopWaitSecs 发送结束进程的信号后等待的秒数
func (that *Entry) StopWaitSecs(defaultVal int) int {
	if that.stopWaitSecs > 0 {
		return that.stopWaitSecs
	}
	return defaultVal
}

// SetStopWaitSecs 发送结束进程的信号后等待的秒数
func (that *Entry) SetStopWaitSecs(stopWaitSecs int) {
	that.stopWaitSecs = stopWaitSecs
}

// KillWaitSecs 强杀进程等待秒数
func (that *Entry) KillWaitSecs(defaultVal int) int {
	if that.killWaitSecs > 0 {
		return that.killWaitSecs
	}
	return defaultVal
}

// SetKillWaitSecs 强杀进程等待秒数
func (that *Entry) SetKillWaitSecs(killWaitSecs int) {
	that.killWaitSecs = killWaitSecs
}

// Environment 获取当前进程的环境变量
func (that *Entry) Environment() []string {
	return that.environment
}

// SetEnvironment 设置当前进程的环境变量
func (that *Entry) SetEnvironment(environment []string) {
	that.environment = environment
}

// Extend 获取扩展配置
func (that *Entry) Extend(key string, defVal string) string {
	if val, found := that.extend[key]; found {
		return val
	}
	return defVal
}

// SetExtend Extend 设置扩展配置信息
func (that *Entry) SetExtend(key string, val string) {
	that.extend[key] = val
}

// GetBytes 通过可识别字符串，返回数字容量
//
//	logSize=1MB
//	logSize=1GB
//	logSize=1KB
//	logSize=1024
//
func (that *Entry) getBytes(value string, defValue int) int {

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

func (that *Entry) toInt(s string, factor int, defValue int) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i * factor
	}
	return defValue
}
