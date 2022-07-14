package dserver

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc"
	"os"
)

type Ctrl struct {
	drpc.CallCtx
}

func (that *Ctrl) Info(_ *string) (*Infos, *drpc.Status) {
	var infos = new(Infos)
	defaultServer.serviceList.Iterator(func(_ string, v interface{}) bool {
		dService := v.(*DService)
		for _, sandbox := range dService.sList.Map() {
			s := sandbox.(*sandboxContainer)
			info := &Info{
				SandBoxName: s.sandbox.Name(),
				ServiceName: dService.Name(),
				Status:      s.state.String(),
				Pid:         gconv.String(os.Getpid()),
				Uptime:      s.started.String(),
			}
			infos.List = append(infos.List, info)
		}
		return true
	})
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
	err := service.stopSandbox(*name)
	if err != nil {
		return nil, drpc.NewStatus(102, err.Error())
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
	err := service.startSandbox(*name)
	if err != nil {
		return nil, drpc.NewStatus(102, err.Error())
	}
	return &Result{}, nil
}
