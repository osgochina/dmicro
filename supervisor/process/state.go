package process

type State int

const (
	// Stopped 已停止
	Stopped State = iota

	// Starting 启动中
	Starting = 10

	// Running 运行中
	Running = 20

	// Backoff 已挂起
	Backoff = 30

	// Stopping 停止中
	Stopping = 40

	// Exited 已退出
	Exited = 100

	// Fatal 启动失败
	Fatal = 200

	// Unknown 未知状态
	Unknown = 1000
)

// String 把进程状态转换成可识别的字符串
func (p State) String() string {
	switch p {
	case Stopped:
		return "Stopped"
	case Starting:
		return "Starting"
	case Running:
		return "Running"
	case Backoff:
		return "Backoff"
	case Stopping:
		return "Stopping"
	case Exited:
		return "Exited"
	case Fatal:
		return "Fatal"
	default:
		return "Unknown"
	}
}

// 更改进程的运行状态
func (that *Process) changeStateTo(procState State) {
	that.state = procState
}
