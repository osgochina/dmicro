package process

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/osgochina/dmicro/logger"
	"sync"
)

type Manager struct {
	processes *gmap.StrAnyMap
}

// NewManager 创建进程管理器
func NewManager() *Manager {
	return &Manager{
		processes: gmap.NewStrAnyMap(true),
	}
}

// NewProcess 创建进程
// path: 可执行文件路径
// args: 参数
// environment: 环境变量
func (that *Manager) NewProcess(path string, args []string, environment []string) *Process {
	p := NewProcess(path, args, environment)
	p.Manager = that
	that.processes.Set(p.GetName(), p)
	logger.Info("创建进程:", p.GetName())
	return p
}

// NewProcessByEntry 创建进程
// entry: 配置对象
func (that *Manager) NewProcessByEntry(entry *Entry) *Process {
	p := NewProcessByEntry(entry)
	p.Manager = that
	that.processes.Set(p.GetName(), p)
	logger.Info("创建进程:", p.GetName())
	return p
}

// NewProcessCmd 创建进程
// path: shell命令
// environment: 环境变量
func (that *Manager) NewProcessCmd(cmd string, environment ...[]string) *Process {
	p := NewProcessCmd(cmd, environment...)
	p.Manager = that
	that.processes.Set(p.GetName(), p)
	logger.Info("创建进程:", p.GetName())
	return p
}

// Add 添加进程到Manager
func (that *Manager) Add(name string, proc *Process) {
	that.processes.Set(name, proc)
	logger.Info("添加进程:", name)
}

// Remove 从Manager移除进程
func (that *Manager) Remove(name string) *Process {
	proc := that.processes.Remove(name)
	if proc == nil {
		return nil
	}
	logger.Info("remove process:", name)
	return proc.(*Process)
}

// Clear 清除进程
func (that *Manager) Clear() {
	that.processes.Clear()
}

// ForEachProcess 迭代进程列表
func (that *Manager) ForEachProcess(procFunc func(p *Process)) {
	that.processes.Iterator(func(k string, v interface{}) bool {
		procFunc(v.(*Process))
		return true
	})
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
	for _, proc := range that.processes.Map() {
		tmpProcList = append(tmpProcList, proc.(*Process))
	}
	return tmpProcList
}

// Find 根据进程名查询进程
func (that *Manager) Find(name string) *Process {
	proc, ok := that.processes.Search(name)
	if ok {
		return proc.(*Process)
	}
	return nil
}

// GetAllProcessInfo 获取所有进程信息
func (that *Manager) GetAllProcessInfo() ([]*Info, error) {
	AllProcessInfo := make([]*Info, 0)
	that.ForEachProcess(func(proc *Process) {
		procInfo := proc.GetProcessInfo()
		AllProcessInfo = append(AllProcessInfo, procInfo)
	})
	return AllProcessInfo, nil
}

// GetProcessInfo 获取指定进程名的进程信息
func (that *Manager) GetProcessInfo(name string) (*Info, error) {
	proc := that.Find(name)
	if proc == nil {
		return nil, fmt.Errorf("no process named %s", name)
	}
	return proc.GetProcessInfo(), nil
}

// StartProcess 启动指定进程
func (that *Manager) StartProcess(name string, wait bool) (bool, error) {
	logger.Infof("启动进程[%s]", name)
	proc := that.Find(name)
	if proc == nil {
		return false, fmt.Errorf("没有找到要启动的进程[%s]", name)
	}
	proc.Start(wait)
	return true, nil
}

// StopProcess 停止指定进程
func (that *Manager) StopProcess(name string, wait bool) (bool, error) {
	logger.Infof("结束进程[%s]", name)
	proc := that.Find(name)
	if proc == nil {
		return false, fmt.Errorf("没有找到要结束的进程[%s]", name)
	}
	proc.Stop(wait)
	return true, nil
}
