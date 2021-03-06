//go:build windows
// +build windows

package graceful

import (
	"github.com/osgochina/dmicro/drpc/netproto/kcp"
	"github.com/osgochina/dmicro/drpc/netproto/quic"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func SetInheritListener(_ []InheritAddr) error { return nil }

func (that *graceful) GraceSignal() {
	// subscribe to SIGINT signals
	signal.Notify(that.signal, os.Interrupt, os.Kill)
	<-that.signal // wait for SIGINT
	signal.Stop(that.signal)
	that.Shutdown()
}
func (that *graceful) Reboot(_ ...time.Duration)                                {}
func (that *graceful) shutdownMaster()                                          {}
func (that *graceful) rebootMasterWorker()                                      {}
func (that *graceful) AddInherited(_ []net.Listener, _ map[string]string)       {}
func (that *graceful) AddInheritedQUIC(_ []*quic.Listener, _ map[string]string) {}
func (that *graceful) AddInheritedKCP(_ []*kcp.Listener, _ map[string]string)   {}
func (that *graceful) startProcess() (*exec.Cmd, error)                         { return nil, nil }
