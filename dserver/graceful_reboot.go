package dserver

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
	"github.com/osgochina/dmicro/supervisor/process"
	"github.com/osgochina/dmicro/utils/signals"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

var originalWD, _ = os.Getwd()

type filer interface {
	File() (*os.File, error)
}

// Reboot 开启优雅的重启流程
func (that *graceful) rebootSingle(timeout ...time.Duration) {
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
	extraFiles := that.getExtraFiles()

	that.inheritedEnv.Set(isChildKey, "true")
	//获取进程启动的原始
	path := os.Args[0]
	err := genv.SetMap(that.inheritedEnv.Map())
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

// GetExtraFiles 获取要继承的fd列表
func (that *graceful) getExtraFiles() []*os.File {
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

	return extraFiles
}

// 获取单进程进程模式下的监听fd列表
func (that *graceful) getEndpointListenerFdMap() map[string]string {
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
		if that.listenAddrMap != nil {
			// 判断监听的是否是http协议。如果是http协议则不返回
			data := that.listenAddrMap.Get(lis.Addr().String())
			if d, ok := data.(InheritAddr); ok && (d.Network == "http" || d.Network == "https") {
				return true
			}
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

// master worker模式重启，就是对子进程发送退出信号
func (that *graceful) rebootMulti() {
	that.dServer.manager.ForEachProcess(func(p *process.Process) {
		pid := p.Pid()
		logger.Printf("主进程:%d 向子进程: %d 发送信号SIGQUIT", os.Getpid(), pid)
		_ = signals.KillPid(pid, signals.ToSignal("SIGQUIT"), false)
	})
	that.processStatus.Set(statusActionNone)
}

type httpListenerFdMap = map[string]string

// 获取监听列表
func (that *graceful) getGHttpListenerFdMap() map[string]httpListenerFdMap {

	if that.dServer.procModel == ProcessModelSingle {
		return nil
	}
	if that.inheritedProcListener.Len() <= 0 {
		return nil
	}
	sfm := make(map[string]httpListenerFdMap)
	m := map[string]string{
		"https": "",
		"http":  "",
	}
	that.inheritedProcListener.Iterator(func(_ int, v interface{}) bool {
		lis, ok := v.(net.Listener)
		if !ok {
			logger.Warningf("inheritedProcListener 不是 net.Listener类型")
			return true
		}
		if that.listenAddrMap == nil {
			return true
		}
		// 判断监听的是否是http协议。如果是http协议则不返回
		data := that.listenAddrMap.Get(lis.Addr().String())

		d, ok := data.(InheritAddr)
		if ok && d.Network != "http" && d.Network != "https" {
			return true
		}
		f, e := lis.(filer).File()
		if e != nil {
			logger.Error(e)
			return true
		}
		str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
		if d.Network == "https" {
			if len(m["https"]) > 0 {
				m["https"] += ","
			}
			m["https"] += str
		} else {
			if len(m["http"]) > 0 {
				m["http"] += ","
			}
			m["http"] += str
		}
		sfm[d.ServerName] = m
		return true
	})
	return sfm
}
