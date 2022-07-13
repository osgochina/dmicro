package dserver

import (
	"github.com/osgochina/dmicro/drpc"
)

type Ctrl struct {
	drpc.CallCtx
}

func (m *Ctrl) Info(arg *[]int) (*Infos, *drpc.Status) {
	var infos = new(Infos)
	infos.List = append(infos.List, &Info{
		Name:   "name1",
		Status: "RUNNING",
		Uptime: "21 day 12 hour",
	})
	infos.List = append(infos.List, &Info{
		Name:   "name2",
		Status: "STOP",
		Uptime: "1111111",
	})
	// response
	return infos, nil
}
