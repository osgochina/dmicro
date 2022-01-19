package easyservice

import (
	"fmt"
)

var defaultService = NewEasyService()

// DefaultService 获取默认的service
func DefaultService() *EasyService {
	return defaultService
}

// GetNextSandBoxId 获取下一个服务沙盒的id
func GetNextSandBoxId() int {
	return sandBoxIdSeq.Add(1)
}

// Setup 启动服务
func Setup(startFunction StartFunc) {
	defaultService.Setup(startFunction)
}

// GetSandBox 获取指定的服务沙盒
func GetSandBox(sandBoxID ...int) ISandBox {
	id := defaultSandBoxId
	if len(sandBoxID) > 0 {
		id = sandBoxID[0]
	}
	return defaultService.GetSandBox(id)
}

// SetSandBox 注册服务沙盒到主服务
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
func Iterator(f func(sandboxId int, sandbox ISandBox)) {
	defaultService.sList.Iterator(func(k int, v interface{}) bool {
		f(k, v.(ISandBox))
		return true
	})
}

// Shutdown 关闭服务
func Shutdown() {
	defaultService.Shutdown()
}

// SetProcessName 设置进程名称
func SetProcessName(name string) {
	//defaultService.setProcessName(name)
}
