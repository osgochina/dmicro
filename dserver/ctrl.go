package dserver

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/drpc/proto/pbproto"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/supervisor/process"
	"os"
	"time"
)

func (that *DServer) endpoint() {
	unix := gfile.TempDir(fmt.Sprintf("%s.sock", that.name))
	if gfile.Exists(unix) {
		_ = gfile.Remove(unix)
	}
	cfg := drpc.EndpointConfig{
		Network:  "unix",
		ListenIP: unix,
	}
	that.ctrlEndpoint = drpc.NewEndpoint(cfg)
	that.ctrlEndpoint.RouteCall(new(Ctrl))
	go func() {
		err := that.ctrlEndpoint.ListenAndServe(pbproto.NewPbProtoFunc())
		if err != nil {
			logger.Warning(err)
		}
		_ = gfile.Remove(unix)
	}()
	time.Sleep(1 * time.Second)
}

type Ctrl struct {
	drpc.CallCtx
}

func (that *Ctrl) Info(_ *string) (*Infos, *drpc.Status) {
	var infos = new(Infos)
	// 单进程
	if defaultServer.procModel == ProcessModelSingle {
		defaultServer.serviceList.Iterator(func(_ interface{}, v interface{}) bool {
			dService := v.(*DService)
			for _, sandbox := range dService.sList.Map() {
				s := sandbox.(*sandboxContainer)
				info := &Info{
					SandBoxName: s.sandbox.Name(),
					ServiceName: dService.Name(),
					Status:      s.state.String(),
					Description: that.createDescription(s.state, s.started, s.stopTime),
				}
				infos.List = append(infos.List, info)
			}
			return true
		})
	}
	// 多进程模式
	if defaultServer.procModel == ProcessModelMulti {
		for _, v := range defaultServer.serviceList.Map() {
			dService := v.(*DService)
			for _, sandbox := range dService.sList.Map() {
				s := sandbox.(*sandboxContainer)
				procInfo, err := defaultServer.manager.GetProcessInfo(dService.Name())
				if err != nil {
					return nil, drpc.NewStatus(100, err.Error())
				}
				info := &Info{
					SandBoxName: s.sandbox.Name(),
					ServiceName: dService.Name(),
					Status:      procInfo.StateName,
					Description: procInfo.Description,
				}
				infos.List = append(infos.List, info)
			}
		}
	}
	return infos, nil
}

// Stop 停止指定的服务
func (that *Ctrl) Stop(name *string) (*Result, *drpc.Status) {
	if len(*name) <= 0 {
		return nil, drpc.NewStatus(100, "未传入sandbox name")
	}
	service, found := defaultServer.searchDServiceBySandboxName(*name)
	if !found {
		return nil, drpc.NewStatus(101, fmt.Sprintf("未找到[%s]", *name))
	}
	// 单进程模式，直接关闭sandbox
	if defaultServer.procModel == ProcessModelSingle {
		err := service.stopSandbox(*name)
		if err != nil {
			return nil, drpc.NewStatus(102, err.Error())
		}
	}
	// 多进程模式，如果关闭sandbox，会把sandbox所在的service全部关闭
	// 暂时不支持关闭单个sandbox功能，后期可以考虑支持
	if defaultServer.procModel == ProcessModelMulti {
		ok, err := defaultServer.manager.StopProcess(service.Name(), true)
		if err != nil {
			return nil, drpc.NewStatus(102, err.Error())
		}
		if !ok {
			return nil, drpc.NewStatus(103, "关闭失败")
		}
	}

	return &Result{}, nil
}

// Start 启动指定的服务
func (that *Ctrl) Start(name *string) (*Result, *drpc.Status) {
	if len(*name) <= 0 {
		return nil, drpc.NewStatus(100, "未传入sandbox name")
	}
	service, found := defaultServer.searchDServiceBySandboxName(*name)
	if !found {
		return nil, drpc.NewStatus(101, fmt.Sprintf("未找到[%s]", *name))
	}
	// 单进程模式，直接开启sandbox
	if defaultServer.procModel == ProcessModelSingle {
		err := service.startSandbox(*name)
		if err != nil {
			return nil, drpc.NewStatus(102, err.Error())
		}
	}
	// 多进程模式，如果启动sandbox，会把sandbox所在的service全部启动
	// 暂时不支持开启单个sandbox功能，后期可以考虑支持
	if defaultServer.procModel == ProcessModelMulti {
		ok, err := defaultServer.manager.StartProcess(service.Name(), true)
		if err != nil {
			return nil, drpc.NewStatus(102, err.Error())
		}
		if !ok {
			return nil, drpc.NewStatus(103, "开启失败")
		}
	}
	return &Result{}, nil
}

// Reload 启动指定的服务
func (that *Ctrl) Reload(name *string) (*Result, *drpc.Status) {
	if len(*name) <= 0 {
		return nil, drpc.NewStatus(100, "未传入sandbox name")
	}
	service, found := defaultServer.searchDServiceBySandboxName(*name)
	if !found {
		return nil, drpc.NewStatus(101, fmt.Sprintf("未找到[%s]", *name))
	}
	// 单进程模式，直接开启sandbox
	if defaultServer.procModel == ProcessModelSingle {
		return nil, drpc.NewStatus(102, "单进程模式不支持reload")
	}
	// 多进程模式，如果启动sandbox，会把sandbox所在的service全部启动
	// 暂时不支持开启单个sandbox功能，后期可以考虑支持
	if defaultServer.procModel == ProcessModelMulti {
		// 判断是否开启了地址继承模式
		if len(defaultServer.inheritAddr) > 0 {
			// 平滑重启逻辑是先启动一个新的继承，再结束老的进程
			ok, err := defaultServer.manager.GracefulReload(service.Name(), true)
			if err != nil {
				return nil, drpc.NewStatus(102, err.Error())
			}
			if !ok {
				return nil, drpc.NewStatus(103, "平滑重启失败")
			}
		} else {
			// 如果未开启地址继承模式，则不需要平滑重启，只需要先结束进程，再启动进程
			ok, err := defaultServer.manager.StopProcess(service.Name(), true)
			if err != nil {
				return nil, drpc.NewStatus(102, err.Error())
			}
			ok, err = defaultServer.manager.StartProcess(service.Name(), true)
			if err != nil {
				return nil, drpc.NewStatus(102, err.Error())
			}
			if !ok {
				return nil, drpc.NewStatus(103, "平滑重启失败")
			}
		}

	}
	return &Result{}, nil
}

// Debug 设置debug模式
func (that *Ctrl) Debug(debug *bool) (*Result, *drpc.Status) {
	if *debug {
		logger.SetDebug(true)
		_ = defaultServer.config.Set("Debug", "true")
	} else {
		logger.SetDebug(false)
		_ = defaultServer.config.Set("Debug", "false")
	}
	return &Result{}, nil
}

// GetDescription 获取进程描述
func (that *Ctrl) createDescription(state process.State, startTime *gtime.Time, stopTime *gtime.Time) string {
	if state == process.Running {
		seconds := int(time.Now().Sub(startTime.Time).Seconds())
		minutes := seconds / 60
		hours := minutes / 60
		days := hours / 24
		if days > 0 {
			return fmt.Sprintf("pid %d, uptime %d days, %d:%02d:%02d", os.Getpid(), days, hours%24, minutes%60, seconds%60)
		}
		return fmt.Sprintf("pid %d, uptime %d:%02d:%02d", os.Getpid(), hours%24, minutes%60, seconds%60)
	} else if state != process.Stopped {
		return stopTime.String()
	}
	return ""
}
