package dserver

import "github.com/gogf/gf/container/gmap"

type DService struct {
	sList *gmap.IntAnyMap //启动的服务列表
}

func (that *DService) AddSandBox(s ISandbox) {
	that.sList.Set(s.ID(), s)
}

// GetSandBox 获取指定的服务沙盒
func (that *DService) GetSandBox(id int) ISandbox {
	s, found := that.sList.Search(id)
	if !found {
		return nil
	}
	return s.(ISandbox)
}
