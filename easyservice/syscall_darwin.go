// +build darwin

package easyservice

import "syscall"

var syscallSIGUSR = syscall.Signal(0)

func syscallKill(_ int, _ syscall.Signal) error {
	return nil
}

// 检查进程是否存在
func checkStart(_ int) bool {
	return false
}

func setProcessName(_ string) {}
