// +build windows

package signals

import (
	"errors"
	"fmt"
	"github.com/osgochina/dmicro/logger"
	"os"
	"os/exec"
	"syscall"
)

// ToSignal 传入信号字符串，返回标准信号
func ToSignal(signalName string) (os.Signal, error) {
	if signalName == "HUP" {
		return syscall.SIGHUP, nil
	} else if signalName == "INT" {
		return syscall.SIGINT, nil
	} else if signalName == "QUIT" {
		return syscall.SIGQUIT, nil
	} else if signalName == "KILL" {
		return syscall.SIGKILL, nil
	} else if signalName == "USR1" {
		logger.Warning("signal USR1 is not supported in windows")
		return nil, errors.New("signal USR1 is not supported in windows")
	} else if signalName == "USR2" {
		logger.Warning("signal USR2 is not supported in windows")
		return nil, errors.New("signal USR2 is not supported in windows")
	} else {
		return syscall.SIGTERM, nil

	}

}

// Kill 向指定的进程发送信号
// process: 进程对象
// sig: 信号
// sigChildren: windows 下会忽略这个参数
func Kill(process *os.Process, sig os.Signal, sigChildren bool) error {
	//Signal command can't kill children processes, call  taskkill command to kill them
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", process.Pid))
	err := cmd.Start()
	if err == nil {
		return cmd.Wait()
	}
	//if fail to find taskkill, fallback to normal signal
	return process.Signal(sig)
}
