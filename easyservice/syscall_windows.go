// +build windows

package easyservice

import "syscall"

var syscallSIGUSR1 = syscall.Signal(0)

func syscallKill(pid int, sig syscall.Signal) (err error) {
	return nil
}

// 检查进程是否存在
func checkStart(pid int) bool {
	return false
}

func setProcessName(name string) {
}
