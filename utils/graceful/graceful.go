package graceful

import (
	"encoding/json"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/util/gconv"
	"os"
	"sync"
	"time"
)

// MinShutdownTimeout 最小停止超时时间
const MinShutdownTimeout = 15 * time.Second

// 当前是否是在子进程
const isWorkerKey = "GRACEFUL_IS_WORKER"

// 父进程的监听列表| Master进程的监听列表
const parentAddrKey = "GRACEFUL_INHERIT_LISTEN_PARENT_ADDR"

// Graceful 优雅重启
type Graceful struct {
	// network:host:[host:port]
	parentAddrList               map[string]map[string][]string
	parentAddrListMutex          sync.Mutex
	locker                       sync.Mutex
	inheritedEnv                 map[string]string
	inheritedProcFiles           []*os.File
	defaultInheritedProcFilesLen int
}

// IsWorker 判断当前进程是在worker进程还是master进程
func (that *Graceful) IsWorker() bool {
	isWorker := genv.GetVar(isWorkerKey, nil)
	if isWorker.IsNil() {
		return false
	}
	if isWorker.Bool() == true {
		return true
	}
	return false
}

// IsMaster 判断当前进程是在worker进程还是master进程
func (that *Graceful) IsMaster() bool {
	return !that.IsWorker()
}

// InitParentAddrList 通过环境变量，初始化父进程监听的端口
//在服务启动的时候，首先从环境变量中获取父进程监听的地址端口，
//如果是首次启动，则不会获取到这些数据
//如果是优雅的无缝重启，则能通过环境变量获取到这些数据，从而复用链接，做到无缝重启
func (that *Graceful) InitParentAddrList() {
	parentAddr := os.Getenv(parentAddrKey)
	_ = json.Unmarshal(gconv.Bytes(parentAddr), &that.parentAddrList)
}

// SetParentAddrList 服务将要重启之前，把当前进程监听的地址端口序列化写入环境变量
func (that *Graceful) SetParentAddrList() {
	b, _ := json.Marshal(that.parentAddrList)
	env := make(map[string]string)
	env[parentAddrKey] = gconv.String(b)
	env[isWorkerKey] = "true"
	that.AddInherited(nil, env)
}

// PushParentAddr 把监听的地址端口写入到变量，优雅结束的时候写入到环境变量，让子进程使用
func (that *Graceful) PushParentAddr(network, host, addr string) {
	that.parentAddrListMutex.Lock()
	defer that.parentAddrListMutex.Unlock()
	that.unifyLocalhost(&host)
	nw, found := that.parentAddrList[network]
	if !found {
		nw = make(map[string][]string)
		that.parentAddrList[network] = nw
	}
	nw[host] = append(nw[host], addr)
}

// PopParentAddr 从监听变量中出栈指定的地址端口
func (that *Graceful) PopParentAddr(network, host, addr string) string {
	that.parentAddrListMutex.Lock()
	defer that.parentAddrListMutex.Unlock()
	that.unifyLocalhost(&host)
	nw, found := that.parentAddrList[network]
	if !found {
		return addr
	}
	h, ok := nw[host]
	if !ok || len(h) == 0 {
		return addr
	}
	nw[host] = h[1:]
	return h[0]
}

// 针对地址格式做统一的转换
func (that *Graceful) unifyLocalhost(host *string) {
	switch *host {
	case "localhost":
		*host = "127.0.0.1"
	case "0.0.0.0":
		*host = "::"
	}
}
