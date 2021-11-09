// +build windows

package graceful

import (
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func (that *graceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, os.Interrupt, os.Kill)
	<-that.signal // wait for SIGINT
	signal.Stop(that.signal)
	that.Shutdown()
}
func (that *graceful) Reboot(timeout ...time.Duration) {}
func (that *graceful) shutdownMaster()                 {}
func (that *graceful) rebootMasterWorker()             {}

func (that *graceful) AddInherited(procFiles []*os.File, envs map[string]string) {}
func SetInheritListener(address []InheritAddr) error {
	return nil
}
func (that *graceful) startProcess() (*exec.Cmd, error) {
	return nil, nil
}

// 发送结束信号给进程
func syscallKillSIGTERM(pid int) error {
	return nil
}
