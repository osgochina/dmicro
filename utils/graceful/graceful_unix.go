// +build linux

package graceful

import (
	"context"
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/inherit"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var originalWD, _ = os.Getwd()

// SetInheritListener 启动master worker模式的监听
func SetInheritListener(address []InheritAddr) error {
	defaultGraceful.SetModel(GraceMasterWorker)
	if !defaultGraceful.isChild() {
		var ch = make(chan int, 1)
		go func() {
			ch <- 1
			defaultGraceful.GraceSignal()
		}()
		<-ch
		for _, addr := range address {
			network := defaultGraceful.translateNetwork(addr.Network)
			if addr.Network == "https" && addr.TlsConfig == nil {
				return gerror.Newf("https 协议，必须传入证书")
			}
			err := defaultGraceful.inheritedListener(inherit.NewFakeAddr(network, addr.Host, addr.Port), addr.TlsConfig)
			if err != nil {
				return err
			}
			defaultGraceful.setMWListenAddr(addr)
			logger.Printf("Master Worker模式，主进程监听(network: %s,host: %s,port: %s)", addr.Network, addr.Host, addr.Port)
		}
		cmd, err := defaultGraceful.startProcess()
		if err != nil {
			return gerror.Newf("启动子进程失败，error:%v", err)
		}
		defaultGraceful.mwChildCmd <- cmd
		defaultGraceful.MWWait()
	}
	return nil
}

// 转换协议
func (that *graceful) translateNetwork(network string) string {
	switch network {
	case "tcp", "tcp4", "tcp6", "http", "https":
		return "tcp"
	case "unix", "unixpacket", "invalid_unix_net_for_test":
		return "unix"
	default:
		return "tcp"
	}
}

// 保存监听的地址
func (that *graceful) setMWListenAddr(addr InheritAddr) {
	if that.mwListenAddr == nil {
		that.mwListenAddr = gmap.NewStrAnyMap(true)
	}
	// 做一次地址转换
	host := addr.Host
	if len(host) == 0 {
		host = "0.0.0.0"
	}
	if host == "0.0.0.0" {
		host = "[::]"
	}
	that.mwListenAddr.Set(fmt.Sprintf("%s:%s", host, addr.Port), addr)
}

// MWWait master worker 模式的主进程等待子进程运行
func (that *graceful) MWWait() {
	for {
		var err error
		select {
		case mwCmd, ok := <-that.mwChildCmd:
			if !ok {
				logger.Fatalf("Master-Worker模式主进程出错")
				return
			}
			that.mwPid = mwCmd.Process.Pid
			logger.Printf("Master-Worker模式启动子进程成功，父进程:%d,子进程:%d", os.Getpid(), that.mwPid)
			err = mwCmd.Wait()
			if err != nil {
				logger.Warningf("子进程:%d 非正常退出，退出原因:%v", that.mwPid, err)
			} else {
				logger.Printf("子进程:%d 正常退出", that.mwPid)
			}
		}
	}
}

// AddInherited 添加需要给重启后新进程继承的文件句柄和环境变量
func (that *graceful) AddInherited(procListener []net.Listener, envs map[string]string) {
	if len(procListener) > 0 {
		for _, f := range procListener {
			// 判断需要添加的文件句柄是否已经存在,不存在才能追加
			if that.inheritedProcListener.Search(f) == -1 {
				that.inheritedProcListener.Append(f)
			}
		}
	}
	if len(envs) > 0 {
		that.inheritedEnv.Sets(envs)
	}
}

// GraceSignal 监听信号
func (that *graceful) GraceSignal() {
	if that.model == GraceChangeProcess {
		that.graceSignalGracefulChangeProcess()
		return
	}
	if that.model == GraceMasterWorker {
		that.graceSignalGracefulMW()
		return
	}
}

// 父子进程模式平滑重启
func (that *graceful) graceSignalGracefulChangeProcess() {
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
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
			signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
			// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
			that.Shutdown(time.Second)
			continue
		// 平滑的关闭服务
		case syscall.SIGTERM:
			signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
			// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
			that.Shutdown()
			continue
		// 平滑重启服务
		case syscall.SIGUSR2:
			that.Reboot()
			continue
		default:
		}
	}
}

// MasterWorker模式平滑重启
func (that *graceful) graceSignalGracefulMW() {
	pid := os.Getpid()
	if that.isChild() {
		signal.Notify(
			that.signal,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGABRT,
		)
		for {
			sig := <-that.signal
			logger.Printf(`进程:%d,收到信号: %s`, pid, sig.String())
			switch sig {
			// 强制关闭服务
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
				that.Shutdown(time.Second)
				continue
			// 平滑的关闭服务
			case syscall.SIGTERM:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
				that.Shutdown()
				continue
			default:
			}
		}
	} else {
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
		for {
			sig := <-that.signal
			logger.Printf(`进程:%d,收到信号: %s`, pid, sig.String())
			switch sig {
			// 关闭服务
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
				that.shutdownMaster()
				continue
			// 平滑重启服务
			case syscall.SIGUSR2:
				that.Reboot()
				continue
			default:
			}
		}
	}
}

// Reboot 开启优雅的重启流程
func (that *graceful) Reboot(timeout ...time.Duration) {
	that.processStatus.Set(statusActionRestarting)
	if that.model == GraceMasterWorker {
		that.rebootMasterWorker()
		return
	}
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

// master worker模式重启，就是对子进程发送退出信号
func (that *graceful) rebootMasterWorker() {
	pid := that.mwPid

	cmd, err := that.startProcess()
	if err != nil {
		logger.Errorf("MasterWorker模式下重启子进程失败,err:%v", err)
		return
	}
	logger.Printf("主进程:%d 向子进程: %d 发送信号SIGTERM", os.Getpid(), pid)
	_ = syscallKillSIGTERM(pid)
	that.processStatus.Set(statusActionNone)
	that.mwChildCmd <- cmd
}

//master worker进程模式，退出master进程方法
func (that *graceful) shutdownMaster() {
	defer os.Exit(0)
	pid := that.mwPid
	logger.Printf(`主进程:%d 向子进程: %d 发送信号SIGTERM`, os.Getpid(), pid)
	_ = syscallKillSIGTERM(pid)
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
	endpointFM := that.getEndpointListenerFdMap()
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

	gHttpSFM := that.getGHttpListenerFdMap()
	if len(gHttpSFM) > 0 {
		for name, m := range gHttpSFM {
			for fdk, fdv := range m {
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
					gHttpSFM[name][fdk] = strings.TrimRight(s, ",")
				}
			}
		}
		buffer, _ := gjson.Encode(gHttpSFM)
		that.inheritedEnv.Set(adminActionReloadEnvKey, string(buffer))
	}

	return extraFiles, nil
}

// 发送结束信号给进程
func syscallKillSIGTERM(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
