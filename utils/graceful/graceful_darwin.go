// +build darwin

package graceful

import (
	"net"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func SetInheritListener(_ []InheritAddr) error { return nil }

// 发送结束信号给进程
func syscallKillSIGTERM(_ int) error { return nil }

func (that *graceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, os.Interrupt, os.Kill)
	<-that.signal // wait for SIGINT
	signal.Stop(that.signal)
	that.Shutdown()
}
func (that *graceful) Reboot(_ ...time.Duration)                             {}
func (that *graceful) shutdownMaster()                                       {}
func (that *graceful) rebootMasterWorker()                                   {}
func (that *graceful) AddInherited(_ []net.Listener, envs map[string]string) {}
func (that *graceful) startProcess() (*exec.Cmd, error)                      { return nil, nil }