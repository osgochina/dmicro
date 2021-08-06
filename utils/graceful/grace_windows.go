// +build windows

package graceful

import (
	"os"
	"os/signal"
)

func (that *Graceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, os.Interrupt, os.Kill)
	<-that.signal // wait for SIGINT
	signal.Stop(that.signal)
	that.Shutdown()
}

func (that *Graceful) Reboot(timeout ...time.Duration) {
	defer os.Exit(0)
	that.logger.Infof("windows system doesn't support reboot! call Shutdown() is recommended.")
}

func (that *Graceful) AddInherited(procFiles []*os.File, envs []*Env) {}
