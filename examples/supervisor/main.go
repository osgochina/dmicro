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
	opts := process.NewProcOptions(
		process.ProcCommand(fmt.Sprintf("%s/../simple/server", gfile.MainPkgPath())),
		process.ProcName("simpleserver"),
		process.ProcUser("lzm"),
		process.ProcDirectory(fmt.Sprintf("%s/../", gfile.MainPkgPath())),
		process.ProcRedirectStderr(true),
		process.ProcStdoutLog("/tmp/tttserver.log", "50M", 10),
	)
	proc, err := m.NewProcessByOptions(opts)
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
