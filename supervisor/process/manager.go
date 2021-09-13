package process

import (
	"fmt"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/supervisor/procconf"
	"sync"
)

type Manager struct {
	procList map[string]*Process
	lock     sync.Mutex
}

// NewManager 创建进程管理器
func NewManager() *Manager {
	return &Manager{
		procList: make(map[string]*Process),
	}
}

// CreateProcess 创建进程
func (that *Manager) CreateProcess(conf *procconf.ProcEntry) *Process {
	that.lock.Lock()
	defer that.lock.Unlock()
	procName := conf.GetProcessName()
	proc, ok := that.procList[procName]
	if !ok {
		proc = NewProcess(conf)
		that.procList[procName] = proc
	}
	logger.Info("create process:", procName)
	return proc
}

// Add 添加进程到Manager
func (that *Manager) Add(name string, proc *Process) {
	that.lock.Lock()
	defer that.lock.Unlock()
	that.procList[name] = proc
	logger.Info("add process:", name)
}

// Remove 从Manager移除进程
func (that *Manager) Remove(name string) *Process {
	that.lock.Lock()
	defer that.lock.Unlock()
	proc, _ := that.procList[name]
	delete(that.procList, name)
	logger.Info("remove process:", name)
	return proc
}

// Clear 清除进程
func (that *Manager) Clear() {
	that.lock.Lock()
	defer that.lock.Unlock()
	that.procList = make(map[string]*Process)
}

// ForEachProcess 迭代进程列表
func (that *Manager) ForEachProcess(procFunc func(p *Process)) {
	that.lock.Lock()
	defer that.lock.Unlock()

	procList := that.getAllProcess()
	for _, proc := range procList {
		procFunc(proc)
	}
}

// StopAllProcesses 关闭所有进程
func (that *Manager) StopAllProcesses() {
	var wg sync.WaitGroup

	that.ForEachProcess(func(proc *Process) {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			proc.Stop(true)
		}(&wg)
	})

	wg.Wait()
}

// 获取所有进程列表
func (that *Manager) getAllProcess() []*Process {
	tmpProcList := make([]*Process, 0)
	for _, proc := range that.procList {
		tmpProcList = append(tmpProcList, proc)
	}
	return tmpProcList
}

// Find 根据进程名查询进程
func (that *Manager) Find(name string) *Process {
	that.lock.Lock()
	defer that.lock.Unlock()
	proc, ok := that.procList[name]
	if ok {
		return proc
	}
	return nil
}

// GetAllProcessInfo 获取所有进程信息
func (that *Manager) GetAllProcessInfo() ([]*ProcessInfo, error) {
	AllProcessInfo := make([]*ProcessInfo, 0)
	that.ForEachProcess(func(proc *Process) {
		procInfo := proc.GetProcessInfo()
		AllProcessInfo = append(AllProcessInfo, procInfo)
	})
	return AllProcessInfo, nil
}

// GetProcessInfo 获取指定进程名的进程信息
func (that *Manager) GetProcessInfo(name string) (*ProcessInfo, error) {
	proc := that.Find(name)
	if proc == nil {
		return nil, fmt.Errorf("no process named %s", name)
	}
	return proc.GetProcessInfo(), nil
}

// StartProcess 启动指定进程
func (that *Manager) StartProcess(name string, wait bool) (bool, error) {
	logger.Infof("start process %s", name)
	proc := that.Find(name)
	if proc == nil {
		return false, fmt.Errorf("fail to find process %s", name)
	}
	proc.Start(wait)
	return true, nil
}

// StopProcess 停止指定进程
func (that *Manager) StopProcess(name string, wait bool) (bool, error) {
	logger.Infof("stop process %s", name)
	proc := that.Find(name)
	if proc == nil {
		return false, fmt.Errorf("fail to find process %s", name)
	}
	proc.Stop(wait)
	return true, nil
}
