package easyservice

import (
	"fmt"
)

//Deprecated
var defaultService = newEasyService()

// DefaultService 获取默认的service
//Deprecated
func DefaultService() *EasyService {
	return defaultService
}

// GetNextSandBoxId 获取下一个服务沙盒的id
//Deprecated
func GetNextSandBoxId() int {
	return sandBoxIdSeq.Add(1)
}

// Setup 启动服务
//Deprecated
func Setup(startFunction StartFunc) {
	defaultService.Setup(startFunction)
}

// GetSandBox 获取指定的服务沙盒
//Deprecated
func GetSandBox(sandBoxID ...int) ISandBox {
	id := defaultSandBoxId
	if len(sandBoxID) > 0 {
		id = sandBoxID[0]
	}
	return defaultService.GetSandBox(id)
}

// SetSandBox 注册服务沙盒到主服务
//Deprecated
func SetSandBox(box ISandBox) {
	if box == nil {
		panic("context: Register backend is nil")
	}
	found := defaultService.GetSandBox(box.ID())
	if found != nil {
		panic(fmt.Sprintf("context: Register called twice for backend %d", box.ID()))
	}
	defaultService.AddSandBox(box)
}

// RemoveSandBox 移除服务沙盒
//Deprecated
func RemoveSandBox(sandBoxID int) error {
	box := defaultService.GetSandBox(sandBoxID)
	err := box.Shutdown()
	if err != nil {
		return err
	}

	defaultService.sList.Remove(sandBoxID)
	return nil
}

// Iterator 迭代服务沙盒
//Deprecated
func Iterator(f func(sandboxId int, sandbox ISandBox)) {
	defaultService.sList.Iterator(func(k int, v interface{}) bool {
		f(k, v.(ISandBox))
		return true
	})
}

// Shutdown 关闭服务
//Deprecated
func Shutdown() {
	defaultService.Shutdown()
}
