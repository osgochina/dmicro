package dserver

var defaultServer = newDServer()

// Setup 启动服务
func Setup(startFunction StartFunc) {
	defaultServer.Setup(startFunction)
}

// Shutdown 关闭服务
func Shutdown() {
	defaultServer.Shutdown()
}
