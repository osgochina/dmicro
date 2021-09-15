// +build !windows,!darwin

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
	"SIGCLD":    syscall.SIGCLD,
	"SIGCONT":   syscall.SIGCONT,
	"SIGFPE":    syscall.SIGFPE,
	"SIGHUP":    syscall.SIGHUP,
	"SIGILL":    syscall.SIGILL,
	"SIGINT":    syscall.SIGINT,
	"SIGIO":     syscall.SIGIO,
	"SIGIOT":    syscall.SIGIOT,
	"SIGKILL":   syscall.SIGKILL,
	"SIGPIPE":   syscall.SIGPIPE,
	"SIGPOLL":   syscall.SIGPOLL,
	"SIGPROF":   syscall.SIGPROF,
	"SIGPWR":    syscall.SIGPWR,
	"SIGQUIT":   syscall.SIGQUIT,
	"SIGSEGV":   syscall.SIGSEGV,
	"SIGSTKFLT": syscall.SIGSTKFLT,
	"SIGSTOP":   syscall.SIGSTOP,
	"SIGSYS":    syscall.SIGSYS,
	"SIGTERM":   syscall.SIGTERM,
	"SIGTRAP":   syscall.SIGTRAP,
	"SIGTSTP":   syscall.SIGTSTP,
	"SIGTTIN":   syscall.SIGTTIN,
	"SIGTTOU":   syscall.SIGTTOU,
	"SIGUNUSED": syscall.SIGUNUSED,
	"SIGURG":    syscall.SIGURG,
	"SIGUSR1":   syscall.SIGUSR1,
	"SIGUSR2":   syscall.SIGUSR2,
	"SIGVTALRM": syscall.SIGVTALRM,
	"SIGWINCH":  syscall.SIGWINCH,
	"SIGXCPU":   syscall.SIGXCPU,
	"SIGXFSZ":   syscall.SIGXFSZ,
}

// ToSignal 传入信号字符串，返回标准信号
func ToSignal(signalName string) (os.Signal, error) {
	if !strings.HasPrefix(signalName, "SIG") {
		signalName = fmt.Sprintf("SIG%s", signalName)
	}
	if sig, ok := signalMap[signalName]; ok {
		return sig, nil
	}
	return syscall.SIGTERM, nil
}

// Kill 向指定的进程发送信号
// process: 进程对象
// sig: 信号
// sigChildren: 如果为true，则信号会发送到该进程的子进程
func Kill(process *os.Process, sig os.Signal, sigChildren bool) error {
	localSig := sig.(syscall.Signal)
	pid := process.Pid
	if sigChildren {
		pid = -pid
	}
	return syscall.Kill(pid, localSig)
}
