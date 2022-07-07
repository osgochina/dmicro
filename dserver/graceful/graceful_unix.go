package graceful

import (
	"context"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/drpc/netproto/kcp"
	"github.com/osgochina/dmicro/drpc/netproto/quic"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/errors"
	"github.com/osgochina/dmicro/utils/signals"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var originalWD, _ = os.Getwd()

type filer interface {
	File() (*os.File, error)
}

// SetInherited 监听的句柄需要被子进程继承，需要设置提取这些句柄，并把它们设置为环境变量
// 该方法留给业务设置，在结束进程的时候调用
func SetInherited(setInherited func() error) {
	defaultGraceful.setInherited = setInherited
}

// AddInherited 添加需要给重启后新进程继承的文件句柄和环境变量
func AddInherited(procListener []net.Listener, envs map[string]string) {
	if len(procListener) > 0 {
		for _, f := range procListener {
			// 判断需要添加的文件句柄是否已经存在,不存在才能追加
			if defaultGraceful.inheritedProcListener.Search(f) == -1 {
				defaultGraceful.inheritedProcListener.Append(f)
			}
		}
	}
	if len(envs) > 0 {
		defaultGraceful.inheritedEnv.Sets(envs)
	}
}

// AddInheritedQUIC 添加quic协议的监听句柄
func AddInheritedQUIC(procListener []*quic.Listener, envs map[string]string) {
	if len(procListener) > 0 {
		for _, f := range procListener {
			// 判断需要添加的文件句柄是否已经存在,不存在才能追加
			if defaultGraceful.inheritedProcListener.Search(f) == -1 {
				defaultGraceful.inheritedProcListener.Append(f)
			}
		}
	}
	if len(envs) > 0 {
		defaultGraceful.inheritedEnv.Sets(envs)
	}
}

// AddInheritedKCP 添加kcp协议的监听句柄
func AddInheritedKCP(procListener []*kcp.Listener, envs map[string]string) {
	if len(procListener) > 0 {
		for _, f := range procListener {
			// 判断需要添加的文件句柄是否已经存在,不存在才能追加
			if defaultGraceful.inheritedProcListener.Search(f) == -1 {
				defaultGraceful.inheritedProcListener.Append(f)
			}
		}
	}
	if len(envs) > 0 {
		defaultGraceful.inheritedEnv.Sets(envs)
	}
}

// 父子进程模型下，当子进程启动成功，发送信号通知父进程
func (that *graceful) onStart() {
	//非子进程，则什么都不走
	if !that.isChild() {
		return
	}
	pPid := syscall.Getppid()
	if pPid != 1 {
		if err := signals.KillPid(pPid, signals.ToSignal("SIGTERM"), false); err != nil {
			logger.Errorf("子进程重启后向父进程发送信号失败，error: %s", err.Error())
			return
		}
		logger.Printf("平滑重启中,子进程[%d]已向父进程[%d]发送信号'SIGTERM'", syscall.Getpid(), pPid)
	}
}

// GraceSignal 监听信号
func (that *graceful) GraceSignal(model int) {
	// 单进程模型
	if model == 0 {
		that.graceSingle()
		return
	}
	// 多进程模型
	//if model == 1 {
	//	that.graceSignalGracefulMW()
	//	return
	//}
}

// SetShutdown 设置退出的基本参数
// timeout 传入负数，表示永远不过期
// firstSweepFunc 收到退出或重启信号后，立刻需要执行的方法
// beforeExitingFunc 处理好收尾动作后，真正要退出了执行的方法
func (that *graceful) SetShutdown(timeout time.Duration, firstSweepFunc, beforeExitingFunc func() error) {
	if timeout < 0 {
		that.shutdownTimeout = 1<<63 - 1
	} else if timeout < minShutdownTimeout {
		that.shutdownTimeout = minShutdownTimeout
	} else {
		that.shutdownTimeout = timeout
	}
	//进程收到退出或重启信号后，需要执行的方法
	if firstSweepFunc == nil {
		firstSweepFunc = func() error { return nil }
	}
	//退出与重启流程最大等待时间
	if beforeExitingFunc == nil {
		beforeExitingFunc = func() error { return nil }
	}
	that.firstSweep = func() error {
		return errors.Merge(
			firstSweepFunc(), //执行自定义方法
			that.setInherited(),
		)
	}
	that.beforeExiting = func() error {
		return errors.Merge(beforeExitingFunc())
	}
}

// 父子进程模式平滑重启
func (that *graceful) graceSingle() {
	signal.Notify(
		that.signal,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	pid := os.Getpid()
	for {
		sig := <-that.signal
		logger.Printf(`进程:%d,收到信号: %s`, pid, sig.String())
		switch sig {
		// 强制关闭服务
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT:
			signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
			// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
			that.shutdown(time.Second)
			continue
		// 平滑的关闭服务
		case syscall.SIGQUIT:
			signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
			// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
			that.shutdown()
			continue
		// 平滑重启服务
		case syscall.SIGUSR2:
			that.reboot()
			continue
		default:
		}
	}
}

// Reboot 开启优雅的重启流程
func (that *graceful) reboot(timeout ...time.Duration) {
	that.processStatus.Set(statusActionRestarting)
	pid := os.Getpid()
	logger.Printf("进程:%d,平滑重启中...", pid)
	that.contextExec(timeout, "reboot", func(_ context.Context) <-chan struct{} {
		endCh := make(chan struct{})

		go func() {
			defer close(endCh)
			if err := that.firstSweep(); err != nil {
				logger.Warningf("进程:%d,平滑重启中 - 执行前置方法失败,error: %s", pid, err.Error())
				os.Exit(-1)
			}

			//启动新的进程
			_, err := that.startProcess()
			// 启动新的进程失败，则表示该进程有问题，直接错误退出
			if err != nil {
				logger.Warningf("进程:%d,平滑重启中 - 启动新的进程失败,error: %s", pid, err.Error())
				os.Exit(-1)
			}
		}()
		return endCh
	})
	logger.Printf("进程:%d,进程已进行平滑重启,等待子进程的信号...", pid)
}

//启动新的进程
func (that *graceful) startProcess() (*exec.Cmd, error) {
	extraFiles, err := that.getExtraFiles()
	if err != nil {
		return nil, err
	}
	that.inheritedEnv.Set(isChildKey, "true")
	//获取进程启动的原始
	path := os.Args[0]
	err = genv.SetMap(that.inheritedEnv.Map())
	if err != nil {
		return nil, err
	}
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	envs := genv.All()
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = extraFiles
	cmd.Env = envs
	cmd.Dir = originalWD
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// 获取要继承的fd列表
func (that *graceful) getExtraFiles() ([]*os.File, error) {
	var extraFiles []*os.File
	endpointFM := that.getEndpointListenerFdMapSingle()
	if len(endpointFM) > 0 {
		for fdk, fdv := range endpointFM {
			if len(fdv) > 0 {
				s := ""
				for _, item := range gstr.SplitAndTrim(fdv, ",") {
					array := strings.Split(item, "#")
					fd := uintptr(gconv.Uint(array[1]))
					if fd > 0 {
						s += fmt.Sprintf("%s#%d,", array[0], 3+len(extraFiles))
						extraFiles = append(extraFiles, os.NewFile(fd, ""))
					} else {
						s += fmt.Sprintf("%s#%d,", array[0], 0)
					}
				}
				endpointFM[fdk] = strings.TrimRight(s, ",")
			}
		}
		buffer, _ := gjson.Encode(endpointFM)
		that.inheritedEnv.Set(parentAddrKey, string(buffer))
	}
	return extraFiles, nil
}

// 获取单进程进程模式下的监听fd列表
func (that *graceful) getEndpointListenerFdMapSingle() map[string]string {
	if that.inheritedProcListener.Len() <= 0 {
		return nil
	}
	m := map[string]string{
		"tcp":  "",
		"quic": "",
		"kcp":  "",
	}
	that.inheritedProcListener.Iterator(func(_ int, v interface{}) bool {
		lis, ok := v.(net.Listener)
		if !ok {
			logger.Warningf("inheritedProcListener 不是 net.Listener类型")
			return true
		}
		quicLis, ok := v.(*quic.Listener)
		if ok {
			f, e := quicLis.PacketConn().(filer).File()
			if e != nil {
				logger.Error(e)
				return false
			}
			str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
			if len(m["quic"]) > 0 {
				m["quic"] += ","
			}
			m["quic"] += str
			return true
		}
		kcpLis, ok := v.(*kcp.Listener)
		if ok {
			f, e := kcpLis.PacketConn().(filer).File()
			if e != nil {
				logger.Error(e)
				return false
			}
			str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
			if len(m["kcp"]) > 0 {
				m["kcp"] += ","
			}
			m["kcp"] += str
			return true
		}
		f, e := lis.(filer).File()
		if e != nil {
			logger.Error(e)
			return false
		}
		str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
		if len(m["tcp"]) > 0 {
			m["tcp"] += ","
		}
		m["tcp"] += str
		return true
	})
	return m
}

// GetInheritedFunc 获取继承的fd
func GetInheritedFunc() []int {
	parentAddr := genv.GetVar(parentAddrKey, nil)
	if parentAddr.IsNil() {
		return nil
	}
	json := gjson.New(parentAddr)
	if json.IsNil() {
		return nil
	}
	var fds []int
	for k, v := range json.Map() {
		// 只能使用tcp协议
		if k != "tcp" {
			continue
		}
		fdv := gconv.String(v)
		if len(fdv) > 0 {
			for _, item := range gstr.SplitAndTrim(fdv, ",") {
				array := strings.Split(item, "#")
				fd := gconv.Int(array[1])
				if fd > 0 {
					fds = append(fds, fd)
				}
			}
		}
	}
	return fds
}

// GetInheritedFuncQUIC 获取继承的fd
func GetInheritedFuncQUIC() []int {
	parentAddr := genv.GetVar(parentAddrKey, nil)
	if parentAddr.IsNil() {
		return nil
	}
	json := gjson.New(parentAddr)
	if json.IsNil() {
		return nil
	}
	var fds []int
	for k, v := range json.Map() {
		// 只能使用quic协议
		if k != "quic" {
			continue
		}
		fdv := gconv.String(v)
		if len(fdv) > 0 {
			for _, item := range gstr.SplitAndTrim(fdv, ",") {
				array := strings.Split(item, "#")
				fd := gconv.Int(array[1])
				if fd > 0 {
					fds = append(fds, fd)
				}
			}
		}
	}
	return fds
}

func GetInheritedFuncKCP() []int {
	parentAddr := genv.GetVar(parentAddrKey, nil)
	if parentAddr.IsNil() {
		return nil
	}
	json := gjson.New(parentAddr)
	if json.IsNil() {
		return nil
	}
	var fds []int
	for k, v := range json.Map() {
		// 只能使用kcp协议
		if k != "kcp" {
			continue
		}
		fdv := gconv.String(v)
		if len(fdv) > 0 {
			for _, item := range gstr.SplitAndTrim(fdv, ",") {
				array := strings.Split(item, "#")
				fd := gconv.Int(array[1])
				if fd > 0 {
					fds = append(fds, fd)
				}
			}
		}
	}
	return fds
}

// isChild 判断当前进程是在子进程还是父进程
func (that *graceful) isChild() bool {
	isWorker := genv.GetVar(isChildKey, nil)
	if isWorker.IsNil() {
		return false
	}
	if isWorker.Bool() {
		return true
	}
	return false
}
