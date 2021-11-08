// +build windows

package graceful

import (
	"os"
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

// 发送结束信号给进程
func SyscallKillSIGTERM(pid int) error {
	return nil
}
