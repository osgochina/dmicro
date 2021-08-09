package easyserver

import "github.com/gogf/gf/container/gtype"

//默认服务沙盒的id
const defaultSandBoxId = 1

// 服务沙盒的自增id
var sandBoxIdSeq = gtype.NewInt(defaultSandBoxId)

type ISandBox interface {
	ID() int          // 沙盒id
	Name() string     // 沙盒名字
	Setup() error     // 启动沙盒
	Shutdown() error  //关闭沙盒
	Service() *Server //返回当前所在服务
}
