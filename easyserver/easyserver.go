package easyserver

import (
	"fmt"
)

var defaultServer = NewServer()

// DefaultServer 获取默认的Server
func DefaultServer() *Server {
	return defaultServer
}

// GetNextSandBoxId 获取下一个服务沙盒的id
func GetNextSandBoxId() int {
	return sandBoxIdSeq.Add(1)
}

// Setup 启动服务
func Setup(startFunction StartFunc) {
	defaultServer.Setup(startFunction)
}

// GetSandBox 获取指定的服务沙盒
func GetSandBox(sandBoxID ...int) ISandBox {
	id := defaultSandBoxId
	if len(sandBoxID) > 0 {
		id = sandBoxID[0]
	}
	return defaultServer.GetSandBox(id)
}

// SetSandBox 注册服务沙盒到主服务
func SetSandBox(box ISandBox) {
	if box == nil {
		panic("context: Register backend is nil")
	}
	found := defaultServer.GetSandBox(box.ID())
	if found != nil {
		panic(fmt.Sprintf("context: Register called twice for backend %d", box.ID()))
	}
	defaultServer.AddSandBox(box)
}

// RemoveSandBox 移除服务沙盒
func RemoveSandBox(sandBoxID int) error {
	box := defaultServer.GetSandBox(sandBoxID)
	err := box.Shutdown()
	if err != nil {
		return err
	}

	defaultServer.sList.Remove(sandBoxID)
	return nil
}

// Shutdown 关闭服务
func Shutdown() {
	defaultServer.Shutdown()
}
