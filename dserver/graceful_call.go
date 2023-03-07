package dserver

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/osgochina/dmicro/drpc/netproto/kcp"
	"github.com/osgochina/dmicro/drpc/netproto/quic"
	"net"
	"strings"
)

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

// AddInheritedQUIC 添加quic协议的监听句柄
func (that *graceful) AddInheritedQUIC(procListener []*quic.Listener, envs map[string]string) {
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

// AddInheritedKCP 添加kcp协议的监听句柄
func (that *graceful) AddInheritedKCP(procListener []*kcp.Listener, envs map[string]string) {
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

// GetInheritedFunc 获取继承的fd
func (that *graceful) GetInheritedFunc() []int {
	parentAddr := genv.Get(parentAddrKey, nil)
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
func (that *graceful) GetInheritedFuncQUIC() []int {
	parentAddr := genv.Get(parentAddrKey, nil)
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

func (that *graceful) GetInheritedFuncKCP() []int {
	parentAddr := genv.Get(parentAddrKey, nil)
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
