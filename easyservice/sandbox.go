package easyservice

import (
	"github.com/gogf/gf/v2/container/gtype"
)

//默认服务沙盒的id
const defaultSandBoxId = 1

// 服务沙盒的自增id
var sandBoxIdSeq = gtype.NewInt(defaultSandBoxId)

// ISandBox 服务沙盒的接口
//Deprecated
type ISandBox interface {
	ID() int               // 沙盒id
	Name() string          // 沙盒名字
	Setup() error          // 启动沙盒
	Shutdown() error       //关闭沙盒
	Service() *EasyService //返回当前所在服务
}
