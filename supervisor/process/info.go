package process

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// GetProcessInfo 获取进程的详情
func (that *Process) GetProcessInfo() *ProcessInfo {
	return &ProcessInfo{
		Name:          that.GetName(),
		Description:   that.GetDescription(),
		Start:         int(that.GetStartTime().Unix()),
		Stop:          int(that.GetStopTime().Unix()),
		Now:           int(time.Now().Unix()),
		State:         int(that.GetState()),
		Statename:     that.GetState().String(),
		Spawnerr:      "",
		Exitstatus:    that.GetExitStatus(),
		Logfile:       that.GetStdoutLogfile(),
		StdoutLogfile: that.GetStdoutLogfile(),
		StderrLogfile: that.GetStderrLogfile(),
		Pid:           that.GetPid()}

}

// GetName 获取进程名
func (that *Process) GetName() string {
	return that.procEntry.GetProcessName()
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
		return that.stopTime.String()
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

// GetPid 获取进程pid，返回0表示进程未启动
func (that *Process) GetPid() int {
	that.lock.RLock()
	defer that.lock.RUnlock()

	if that.state == Stopped || that.state == Fatal || that.state == Unknown || that.state == Exited || that.state == Backoff {
		return 0
	}
	return that.cmd.Process.Pid
}

// GetStdoutLogfile 获取标准输出将要写入的日志文件
func (that *Process) GetStdoutLogfile() string {
	fileName := that.procEntry.GetStdoutLogfile("/dev/null")
	expandFile, err := PathExpand(fileName)
	if err != nil {
		return fileName
	}
	return expandFile
}

// GetStderrLogfile 获取标准错误将要写入的日志文件
func (that *Process) GetStderrLogfile() string {
	fileName := that.procEntry.GetStderrLogfile("/dev/null")
	expandFile, err := PathExpand(fileName)
	if err != nil {
		return fileName
	}
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
	strExitCodes := strings.Split(that.procEntry.GetExitCodes("0,2"), ",")
	result := make([]int, 0)
	for _, val := range strExitCodes {
		i, err := strconv.Atoi(val)
		if err == nil {
			result = append(result, i)
		}
	}
	return result
}
