package dserver

import "github.com/gogf/gf/container/gtype"

//默认服务沙盒的id
const defaultSandboxId = 1

// 服务沙盒的自增id
var sandBoxIdSeq = gtype.NewInt(defaultSandboxId)

// ISandbox 服务沙盒的接口
type ISandbox interface {
	ID() int            // 沙盒id
	Name() string       // 沙盒名字
	Setup() error       // 启动沙盒
	Shutdown() error    //关闭沙盒
	Service() *DService //返回当前所在服务
}

// GetNextSandBoxId 获取下一个服务沙盒的id
func GetNextSandBoxId() int {
	return sandBoxIdSeq.Add(1)
}
