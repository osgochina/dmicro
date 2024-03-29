package process

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/osgochina/dmicro/logger"
	"sync"
)

type Manager struct {
	processes *gmap.StrAnyMap
}

// NewManager 创建进程管理器
func NewManager() *Manager {
	return &Manager{
		processes: gmap.NewStrAnyMap(false),
	}
}

// NewProcess 创建进程
// path: 可执行文件路径
// args: 参数
// environment: 环境变量
func (that *Manager) NewProcess(path string, args []string, environment map[string]string, opts ...ProcOption) (*Process, error) {

	opts = append(opts,
		ProcCommand(path),
		ProcArgs(args...),
		ProcEnvironment(environment),
	)
	p := NewProcess(opts...)
	if _, found := that.processes.Search(p.GetName()); found {
		return nil, gerror.Newf("进程[%s]已存在", p.GetName())
	}
	p.Manager = that
	that.processes.Set(p.GetName(), p)
	logger.Info(context.TODO(), "创建进程:", p.GetName())
	return p, nil
}

// NewProcessByOptions 创建进程
// entry: 配置对象
func (that *Manager) NewProcessByOptions(opts ProcOptions) (*Process, error) {
	p := NewProcessByOptions(opts)
	if _, found := that.processes.Search(p.GetName()); found {
		return nil, gerror.Newf("进程[%s]已存在", p.GetName())
	}
	p.Manager = that
	that.processes.Set(p.GetName(), p)
	logger.Info(context.TODO(), "创建进程:", p.GetName())
	return p, nil
}

// NewProcessByProcess 创建进程
// proc: Process对象
func (that *Manager) NewProcessByProcess(proc *Process) (*Process, error) {
	if _, found := that.processes.Search(proc.GetName()); found {
		return nil, gerror.Newf("进程[%s]已存在", proc.GetName())
	}
	proc.Manager = that
	that.processes.Set(proc.GetName(), proc)
	logger.Info(context.TODO(), "创建进程:", proc.GetName())
	return proc, nil
}

// NewProcessCmd 创建进程
// path: shell命令
// environment: 环境变量
func (that *Manager) NewProcessCmd(cmd string, environment map[string]string) (*Process, error) {
	p := NewProcessCmd(cmd, environment)
	if _, found := that.processes.Search(p.GetName()); found {
		return nil, gerror.Newf("进程[%s]已存在", p.GetName())
	}
	p.Manager = that
	that.processes.Set(p.GetName(), p)
	logger.Info(context.TODO(), "创建进程:", p.GetName())
	return p, nil
}

// Add 添加进程到Manager
func (that *Manager) Add(name string, proc *Process) {
	that.processes.Set(name, proc)
	logger.Info(context.TODO(), "添加进程:", name)
}

// Remove 从Manager移除进程
func (that *Manager) Remove(name string) *Process {
	proc := that.processes.Remove(name)
	if proc == nil {
		return nil
	}
	logger.Info(context.TODO(), "remove process:", name)
	return proc.(*Process)
}

// Clear 清除进程
func (that *Manager) Clear() {
	that.processes.Clear()
}

// ForEachProcess 迭代进程列表
func (that *Manager) ForEachProcess(procFunc func(p *Process)) {
	that.processes.Iterator(func(_ string, v interface{}) bool {
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
	logger.Infof(context.TODO(), "启动进程[%s]", name)
	proc := that.Find(name)
	if proc == nil {
		return false, fmt.Errorf("没有找到要启动的进程[%s]", name)
	}
	proc.Start(wait)
	return true, nil
}

// StopProcess 停止指定进程
func (that *Manager) StopProcess(name string, wait bool) (bool, error) {
	logger.Infof(context.TODO(), "结束进程[%s]", name)
	proc := that.Find(name)
	if proc == nil {
		return false, fmt.Errorf("没有找到要结束的进程[%s]", name)
	}
	proc.Stop(wait)
	return true, nil
}

// GracefulReload 停止指定进程
func (that *Manager) GracefulReload(name string, wait bool) (bool, error) {
	logger.Infof(context.TODO(), "平滑重启进程[%s]", name)
	proc := that.Find(name)
	if proc == nil {
		return false, fmt.Errorf("没有找到要重启的进程[%s]", name)
	}
	procClone, err := proc.Clone()
	if err != nil {
		return false, err
	}
	procClone.Start(wait)
	proc.Stop(wait)
	that.processes.Set(name, procClone)
	return true, nil
}
