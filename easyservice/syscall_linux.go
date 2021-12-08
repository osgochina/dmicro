// +build linux

package easyservice

import (
	"syscall"
)

var syscallSIGUSR = syscall.SIGUSR2

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
func setProcessName(_ string) {
	// TODO 该依赖库有点问题，在golang:alpine中无法引入stdlib.h，暂时不支持，后续想到办法
	//gspt.SetProcTitle(name)
}
