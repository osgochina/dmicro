//go:build windows
// +build windows

package dserver

import (
	"os"
	"os/signal"
)

func (that *graceful) graceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, os.Interrupt, os.Kill)
	<-that.signal // wait for SIGINT
	signal.Stop(that.signal)
	that.shutdownSingle()
}
