package main

import (
	loggerv2 "github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/supervisor/procconf"
	"github.com/osgochina/dmicro/supervisor/process"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	process.ReapZombie()
	runServer()
}

func runServer() {
	m := process.NewManager()
	entry := procconf.NewProcEntry()
	entry.SetProcessName("server")
	entry.User = "lzm"
	entry.Directory = "/home/lzm/go/src/dmicro/.examples/"
	entry.Command = "/home/lzm/go/src/dmicro/.examples/simple/server"
	entry.SetRedirectStderr(true)
	entry.SetStdoutLogfile("/tmp/tttserver.log")
	proc := m.CreateProcess(entry)
	proc.Start(false)
	initSignals(m)
	select {}
}

func initSignals(s *process.Manager) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		loggerv2.Infof("receive a signal %s to stop all process & exit", sig)
		s.StopAllProcesses()
		os.Exit(-1)
	}()

}
