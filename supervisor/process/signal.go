package process

import (
	"fmt"
	loggerv2 "github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/supervisor/signals"
	"os"
)

// Signal sends signal to the process
//
// Args:
//   sig - the signal to the process
//   sigChildren - if true, sends the same signal to the process and its children
//
func (that *Process) Signal(sig os.Signal, sigChildren bool) error {
	that.lock.RLock()
	defer that.lock.RUnlock()

	return that.sendSignal(sig, sigChildren)
}

func (that *Process) sendSignals(sigs []string, sigChildren bool) {
	that.lock.RLock()
	defer that.lock.RUnlock()

	for _, strSig := range sigs {
		sig, err := signals.ToSignal(strSig)
		if err == nil {
			_ = that.sendSignal(sig, sigChildren)
		} else {
			loggerv2.Info("program %s,Invalid signal name %s", that.GetName(), strSig)
		}
	}
}

// send signal to the process
//
// Args:
//    sig - the signal to be sent
//    sigChildren - if true, the signal also will be sent to children processes too
//
func (that *Process) sendSignal(sig os.Signal, sigChildren bool) error {
	if that.cmd != nil && that.cmd.Process != nil {
		loggerv2.Infof("Send signal %s to program %s", sig, that.GetName())
		err := signals.Kill(that.cmd.Process, sig, sigChildren)
		return err
	}
	return fmt.Errorf("process is not started")
}
