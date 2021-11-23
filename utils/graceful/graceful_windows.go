// +build windows

package graceful

import (
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func SetInheritListener(address []InheritAddr) error { return nil }

// 发送结束信号给进程
func syscallKillSIGTERM(pid int) error { return nil }

func (that *graceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, os.Interrupt, os.Kill)
	<-that.signal // wait for SIGINT
	signal.Stop(that.signal)
	that.Shutdown()
}
func (that *graceful) Reboot(timeout ...time.Duration)                                  {}
func (that *graceful) shutdownMaster()                                                  {}
func (that *graceful) rebootMasterWorker()                                              {}
func (that *graceful) AddInherited(procListener []net.Listener, envs map[string]string) {}
func (that *graceful) startProcess() (*exec.Cmd, error)                                 { return nil, nil }
