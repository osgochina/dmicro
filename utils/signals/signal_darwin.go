// +build darwin

package signals

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

// 可识别的信号列表
var signalMap = map[string]os.Signal{
	"SIGABRT":   syscall.SIGABRT,
	"SIGALRM":   syscall.SIGALRM,
	"SIGBUS":    syscall.SIGBUS,
	"SIGCHLD":   syscall.SIGCHLD,
	"SIGCONT":   syscall.SIGCONT,
	"SIGEMT":    syscall.SIGEMT,
	"SIGFPE":    syscall.SIGFPE,
	"SIGHUP":    syscall.SIGHUP,
	"SIGILL":    syscall.SIGILL,
	"SIGINFO":   syscall.SIGINFO,
	"SIGINT":    syscall.SIGINT,
	"SIGIO":     syscall.SIGIO,
	"SIGIOT":    syscall.SIGIOT,
	"SIGKILL":   syscall.SIGKILL,
	"SIGPIPE":   syscall.SIGPIPE,
	"SIGPROF":   syscall.SIGPROF,
	"SIGQUIT":   syscall.SIGQUIT,
	"SIGSEGV":   syscall.SIGSEGV,
	"SIGSTOP":   syscall.SIGSTOP,
	"SIGSYS":    syscall.SIGSYS,
	"SIGTERM":   syscall.SIGTERM,
	"SIGTRAP":   syscall.SIGTRAP,
	"SIGTSTP":   syscall.SIGTSTP,
	"SIGTTIN":   syscall.SIGTTIN,
	"SIGTTOU":   syscall.SIGTTOU,
	"SIGURG":    syscall.SIGURG,
	"SIGUSR1":   syscall.SIGUSR1,
	"SIGUSR2":   syscall.SIGUSR2,
	"SIGVTALRM": syscall.SIGVTALRM,
	"SIGWINCH":  syscall.SIGWINCH,
	"SIGXCPU":   syscall.SIGXCPU,
	"SIGXFSZ":   syscall.SIGXFSZ,
}

// ToSignal 传入信号字符串，返回标准信号
func ToSignal(signalName string) os.Signal {
	if !strings.HasPrefix(signalName, "SIG") {
		signalName = fmt.Sprintf("SIG%s", signalName)
	}
	if sig, ok := signalMap[signalName]; ok {
		return sig
	}
	return syscall.SIGTERM

}

// Kill 向指定的进程发送信号
// process: 进程对象
// sig: 信号
// sigChildren: 如果为true，则信号会发送到该进程的子进程
func Kill(process *os.Process, sig os.Signal, sigChildren ...bool) error {
	localSig := sig.(syscall.Signal)
	pid := process.Pid
	if len(sigChildren) > 0 && sigChildren[0] {
		pid = -pid
	}
	return syscall.Kill(pid, localSig)
}

// KillPid 向指定的Pid发送信号
// pid: 进程pid
// sig: 信号
// sigChildren: 如果为true，则信号会发送到该进程的子进程
func KillPid(pid int, sig os.Signal, sigChildren ...bool) error {
	localSig := sig.(syscall.Signal)
	if len(sigChildren) > 0 && sigChildren[0] {
		pid = -pid
	}
	return syscall.Kill(pid, localSig)
}

// 检查进程是否存在
func CheckPidExist(_ int) bool {
	return false
}
