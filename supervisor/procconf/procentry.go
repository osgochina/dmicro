package procconf

import (
	"github.com/gogf/gf/util/gconv"
	"os/exec"
	"strconv"
	"syscall"
)

type ProcEntry struct {
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
	stopAsGroup bool
	//默认为false，向进程组发送kill信号，包括子进程
	killAsGroup bool
	//结束进程发送的信号
	stopSignal string
	// 发送结束进程的信号后等待的秒数
	stopWaitSecs int
	//
	killWaitSecs int
	// 环境变量
	environment map[string]string
	//当进程的二进制文件有修改，是否需要重启
	//RestartWhenBinaryChanged bool
	// 扩展配置
	extend map[string]string
}

func NewProcEntry(command string, args ...string) *ProcEntry {
	proc := &ProcEntry{
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
		stdoutLogFileMaxBytes: "50MB",
		stdoutLogFileBackups:  10,
		redirectStderr:        false,
		stderrLogfile:         "",
		stderrLogFileMaxBytes: "50MB",
		stderrLogFileBackups:  10,
		stopAsGroup:           false,
		killAsGroup:           false,
		environment:           make(map[string]string),
		extend:                make(map[string]string),
	}
	proc.args = append(proc.args, args...)
	return proc
}

// Name 获取进程名称
func (that *ProcEntry) Name() string {
	return that.name
}

// SetName 设置进程名
func (that *ProcEntry) SetName(name string) {
	that.name = name
}

// Command 获取命令
func (that *ProcEntry) Command() string {
	return that.command
}

// SetCommand 设置命令
func (that *ProcEntry) SetCommand(command string) {
	that.command = command
}

func (that *ProcEntry) Args() []string {
	return that.args
}

func (that *ProcEntry) SetArgs(args []string) {
	that.args = args
}

// CreateCommand 生成命令
func (that *ProcEntry) CreateCommand() (*exec.Cmd, error) {

	cmd := exec.Command(that.command)
	if len(that.args) > 0 {
		cmd.Args = append([]string{that.command}, that.args...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	return cmd, nil
}

func (that *ProcEntry) Directory() string {
	return that.directory
}

func (that *ProcEntry) SetDirectory(directory string) {
	that.directory = directory
}

func (that *ProcEntry) AutoStart() bool {
	return that.autoStart
}

func (that *ProcEntry) SetAutoStart(autoStart bool) {
	that.autoStart = autoStart
}

func (that *ProcEntry) StartSecs() int {
	return that.startSecs
}

func (that *ProcEntry) SetStartSecs(startSecs int) {
	that.startSecs = startSecs
}

// AutoReStart 值为 true，false，unexpected
func (that *ProcEntry) AutoReStart(defVal string) string {
	if len(that.autoReStart) > 0 {
		return that.autoReStart
	}
	return defVal
}

// SetAutoReStart 设置自动重启 值为 true，false，unexpected
func (that *ProcEntry) SetAutoReStart(val string) {
	that.autoReStart = val
}

func (that *ProcEntry) ExitCodes(defaultVal string) string {
	if len(that.exitCodes) > 0 {
		return that.exitCodes
	}
	return defaultVal
}

func (that *ProcEntry) SetExitCodes(val string) {
	that.exitCodes = val
}

func (that *ProcEntry) StartRetries() int {
	return that.startRetries
}

func (that *ProcEntry) SetStartRetries(startRetries int) {
	that.startRetries = startRetries
}

func (that *ProcEntry) RestartPause() int {
	return that.restartPause
}

func (that *ProcEntry) SetRestartPause(restartPause int) {
	that.restartPause = restartPause
}

func (that *ProcEntry) User() string {
	return that.user
}

func (that *ProcEntry) SetUser(user string) {
	that.user = user
}

func (that *ProcEntry) Priority() int {
	return that.priority
}

func (that *ProcEntry) SetPriority(priority int) {
	that.priority = priority
}

func (that *ProcEntry) StdoutLogfile(defaultVal string) string {
	if len(that.stdoutLogfile) > 0 {
		return that.stdoutLogfile
	}
	return defaultVal
}

func (that *ProcEntry) SetStdoutLogfile(v string) {
	that.stdoutLogfile = v
}

func (that *ProcEntry) StdoutLogFileMaxBytes(defaultVal int) int {
	return that.getBytes(that.stdoutLogFileMaxBytes, defaultVal)
}

func (that *ProcEntry) SetStdoutLogFileMaxBytes(v string) {
	that.stdoutLogFileMaxBytes = v
}

func (that *ProcEntry) StdoutLogFileBackups(defaultVal int) int {
	if that.stdoutLogFileBackups > 0 {
		return that.stdoutLogFileBackups
	}
	return defaultVal
}

func (that *ProcEntry) SetStdoutLogFileBackups(v int) {
	that.stdoutLogFileBackups = v
}

func (that *ProcEntry) RedirectStderr() bool {
	return that.redirectStderr
}

func (that *ProcEntry) SetRedirectStderr(v bool) {
	that.redirectStderr = v
}

func (that *ProcEntry) StderrLogfile(defaultVal string) string {
	if len(that.stderrLogfile) > 0 {
		return that.stderrLogfile
	}
	return defaultVal
}

func (that *ProcEntry) SetStderrLogfile(v string) {
	that.stderrLogfile = v
}

func (that *ProcEntry) StderrLogFileMaxBytes(defaultVal int) int {
	return that.getBytes(that.stderrLogFileMaxBytes, defaultVal)
}

func (that *ProcEntry) SetStderrLogFileMaxBytes(v string) {
	that.stderrLogFileMaxBytes = v
}

func (that *ProcEntry) StderrLogFileBackups(defaultVal int) int {
	if that.stderrLogFileBackups > 0 {
		return that.stderrLogFileBackups
	}
	return defaultVal
}

func (that *ProcEntry) SetStderrLogFileBackups(v int) {
	that.stderrLogFileBackups = v
}

func (that *ProcEntry) StopAsGroup() bool {
	return that.stopAsGroup
}

func (that *ProcEntry) SetStopAsGroup(stopAsGroup bool) {
	that.stopAsGroup = stopAsGroup
}

func (that *ProcEntry) KillAsGroup() bool {
	return that.killAsGroup
}

func (that *ProcEntry) SetKillAsGroup(killAsGroup bool) {
	that.killAsGroup = killAsGroup
}

func (that *ProcEntry) StopSignal() string {
	return that.stopSignal
}

func (that *ProcEntry) SetStopSignal(stopSignal string) {
	that.stopSignal = stopSignal
}

func (that *ProcEntry) StopWaitSecs(defaultVal int) int {
	if that.stopWaitSecs > 0 {
		return that.stopWaitSecs
	}
	return defaultVal
}

func (that *ProcEntry) SetStopWaitSecs(stopWaitSecs int) {
	that.stopWaitSecs = stopWaitSecs
}

func (that *ProcEntry) KillWaitSecs(defaultVal int) int {
	if that.killWaitSecs > 0 {
		return that.killWaitSecs
	}
	return defaultVal
}
func (that *ProcEntry) SetKillWaitSecs(killWaitSecs int) {
	that.killWaitSecs = killWaitSecs
}

func (that *ProcEntry) Environment() map[string]string {
	return that.environment
}

func (that *ProcEntry) SetEnvironment(environment map[string]string) {
	that.environment = environment
}

func (that *ProcEntry) Extend() map[string]string {
	return that.extend
}

func (that *ProcEntry) SetExtend(extend map[string]string) {
	that.extend = extend
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
