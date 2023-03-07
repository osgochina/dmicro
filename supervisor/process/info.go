package process

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"syscall"
	"time"
)

// Info 进程的运行状态
type Info struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Start         int    `json:"start"`
	Stop          int    `json:"stop"`
	Now           int    `json:"now"`
	State         int    `json:"state"`
	StateName     string `json:"statename"`
	SpawnErr      string `json:"spawnerr"`
	ExitStatus    int    `json:"exitstatus"`
	Logfile       string `json:"logfile"`
	StdoutLogfile string `json:"stdout_logfile"`
	StderrLogfile string `json:"stderr_logfile"`
	Pid           int    `json:"pid"`
}

// GetProcessInfo 获取进程的详情
func (that *Process) GetProcessInfo() *Info {
	return &Info{
		Name:          that.GetName(),
		Description:   that.GetDescription(),
		Start:         int(that.GetStartTime().Unix()),
		Stop:          int(that.GetStopTime().Unix()),
		Now:           int(time.Now().Unix()),
		State:         int(that.GetState()),
		StateName:     that.GetState().String(),
		SpawnErr:      "",
		ExitStatus:    that.GetExitStatus(),
		Logfile:       that.GetStdoutLogfile(),
		StdoutLogfile: that.GetStdoutLogfile(),
		StderrLogfile: that.GetStderrLogfile(),
		Pid:           that.Pid()}

}

// GetName 获取进程名
func (that *Process) GetName() string {
	return that.option.Name
}

// GetDescription 获取进程描述
func (that *Process) GetDescription() string {
	that.lock.RLock()
	defer that.lock.RUnlock()
	if that.state == Running {
		seconds := int(time.Now().Sub(that.startTime).Seconds())
		minutes := seconds / 60
		hours := minutes / 60
		days := hours / 24
		if days > 0 {
			return fmt.Sprintf("pid %d, uptime %d days, %d:%02d:%02d", that.cmd.Process.Pid, days, hours%24, minutes%60, seconds%60)
		}
		return fmt.Sprintf("pid %d, uptime %d:%02d:%02d", that.cmd.Process.Pid, hours%24, minutes%60, seconds%60)
	} else if that.state != Stopped {
		return gtime.New(that.stopTime).String()
	}
	return ""
}

// GetState 获取进程状态
func (that *Process) GetState() State {
	return that.state
}

// GetStartTime 获取进程启动时间
func (that *Process) GetStartTime() time.Time {
	return that.startTime
}

// GetStopTime 获取进程结束时间
func (that *Process) GetStopTime() time.Time {
	switch that.state {
	case Starting:
		fallthrough
	case Running:
		fallthrough
	case Stopping:
		return time.Unix(0, 0)
	default:
		return that.stopTime
	}
}

// GetExitStatus 获取进程退出状态
func (that *Process) GetExitStatus() int {
	that.lock.RLock()
	defer that.lock.RUnlock()

	if that.state == Exited || that.state == Backoff {
		if that.cmd.ProcessState == nil {
			return 0
		}
		status, ok := that.cmd.ProcessState.Sys().(syscall.WaitStatus)
		if ok {
			return status.ExitStatus()
		}
	}
	return 0
}

// Pid 获取进程pid，返回0表示进程未启动
func (that *Process) Pid() int {

	if that.state == Stopped || that.state == Fatal || that.state == Unknown || that.state == Exited || that.state == Backoff {
		return 0
	}
	return that.cmd.Process.Pid
}

// GetStdoutLogfile 获取标准输出将要写入的日志文件
func (that *Process) GetStdoutLogfile() string {
	fileName := "/dev/null"
	if len(that.option.StdoutLogfile) > 0 {
		fileName = that.option.StdoutLogfile
	}
	expandFile := gfile.RealPath(fileName)
	return expandFile
}

// GetStderrLogfile 获取标准错误将要写入的日志文件
func (that *Process) GetStderrLogfile() string {
	fileName := "/dev/null"
	if len(that.option.StderrLogfile) > 0 {
		fileName = that.option.StdoutLogfile
	}
	expandFile := gfile.RealPath(fileName)
	return expandFile
}

// GetStatus 获取进程当前状态
func (that *Process) GetStatus() string {
	if that.cmd.ProcessState.Exited() {
		return that.cmd.ProcessState.String()
	}
	return "running"
}

// 获取进程的退出code值
func (that *Process) getExitCode() (int, error) {
	if that.cmd.ProcessState == nil {
		return -1, fmt.Errorf("no exit code")
	}
	if status, ok := that.cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		return status.ExitStatus(), nil
	}

	return -1, fmt.Errorf("no exit code")

}

// 进程的退出code值是否在设置中的codes列表中
func (that *Process) inExitCodes(exitCode int) bool {
	for _, code := range that.getExitCodes() {
		if code == exitCode {
			return true
		}
	}
	return false
}

// 获取配置的退出code值列表
func (that *Process) getExitCodes() []int {
	strExitCodes := that.option.ExitCodes
	if len(that.option.ExitCodes) > 0 {
		return strExitCodes
	}
	return []int{0, 2}
}
