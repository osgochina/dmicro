// +build linux

package easyservice

import (
	"github.com/osgochina/dmicro/utils/gspt"
	"syscall"
)

var syscallSIGUSR2 = syscall.SIGUSR2

func syscallKill(pid int, sig syscall.Signal) (err error) {
	return syscall.Kill(pid, sig)
}

// 检查进程是否存在
func checkStart(pid int) bool {
	err := syscallKill(pid, 0)
	if err != nil && err.Error() != "no such process" {
		return true
	}
	return false
}

//设置进程名
func setProcessName(name string) {
	gspt.SetProcTitle(name)
}
