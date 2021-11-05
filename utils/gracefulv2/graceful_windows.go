// +build windows

package gracefulv2

import (
	"os"
)

func (that *ChangeProcessGraceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, os.Interrupt, os.Kill)
	<-that.signal // wait for SIGINT
	signal.Stop(that.signal)
	that.Shutdown()
}

func (that *ChangeProcessGraceful) AddInherited(procFiles []*os.File, envs map[string]string) {}

// 发送结束信号给进程
func SyscallKillSIGTERM(pid int) error {
	return nil
}
