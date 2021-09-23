package main

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/osgochina/dmicro/logger"
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
	entry := process.NewEntry(fmt.Sprintf("%s/../simple/server", gfile.MainPkgPath()))
	entry.SetName("simpleserver")
	entry.SetUser("lzm")
	entry.SetDirectory(fmt.Sprintf("%s/../", gfile.MainPkgPath()))
	entry.SetRedirectStderr(true)
	entry.SetStdoutLogfile("/tmp/tttserver.log")
	proc, err := m.NewProcessByEntry(entry)
	if err != nil {
		logger.Error(err)
		return
	}
	proc.Start(true)
	initSignals(m)
	select {}
}

func initSignals(s *process.Manager) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Infof("receive a signal %s to stop all process & exit", sig)
		s.StopAllProcesses()
		os.Exit(-1)
	}()

}
