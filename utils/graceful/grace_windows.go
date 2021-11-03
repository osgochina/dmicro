// +build windows

package graceful

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

var isReboot = false

func (that *Graceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, os.Interrupt, os.Kill)
	<-that.signal // wait for SIGINT
	signal.Stop(that.signal)
	that.Shutdown()
}

func (that *Graceful) Reboot(timeout ...time.Duration) {
	defer os.Exit(0)
	fmt.Println("windows system doesn't support reboot! call Shutdown() is recommended.")
}

func (that *Graceful) AddInherited(procFiles []*os.File, envs map[string]string) {}

// 发送结束信号给进程
func SyscallKillSIGTERM(pid int) error {
	return nil
}
