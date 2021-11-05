package drpc

import (
	"github.com/gogf/gf/container/gset"
	"github.com/osgochina/dmicro/utils/errors"
	"github.com/osgochina/dmicro/utils/gracefulv2"
)

func init() {
	gracefulv2.GetGraceful().SetShutdownEndpoint(func(endpointList *gset.Set) error {
		var count int
		var errCh = make(chan error, endpointList.Size())
		//异步关闭端点
		for _, val := range endpointList.Slice() {
			count++
			e := val.(*endpoint)
			go func(e *endpoint) {
				errCh <- e.Close()
			}(e)
		}
		var err error
		for i := 0; i < count; i++ {
			err = errors.Merge(err, <-errCh)
		}
		close(errCh)
		return err
	})
}

//var drpcGraceful *graceful.ChangeProcessGraceful
//
//func init() {
//	drpcGraceful = graceful.NewChangeProcessGraceful()
//	drpcGraceful.InitParentAddrList()
//	inherit.AddInheritedFunc(drpcGraceful.AddInherited)
//	SetShutdown(5*time.Second, nil, nil)
//}
//
//func GraceSignal() {
//	drpcGraceful.GraceSignal()
//}
//
//// GraceOnStart 启动成功，发送信号
//func GraceOnStart() {
//	onServeListener(nil)
//}

//var endpointList = struct {
//	list map[*endpoint]struct{}
//	rwMu sync.RWMutex
//}{
//	list: make(map[*endpoint]struct{}),
//}
//
//// 新增一个端点
//func addEndpoint(e *endpoint) {
//	endpointList.rwMu.Lock()
//	endpointList.list[e] = struct{}{}
//	endpointList.rwMu.Unlock()
//}
//
////删除一个端点
//func deleteEndpoint(e *endpoint) {
//	endpointList.rwMu.Lock()
//	delete(endpointList.list, e)
//	endpointList.rwMu.Unlock()
//}

// 服务监听成功的情况下，调用该方法。
//如果是在平滑重启的子进程中，则子进程会发送信号给父进程，让父进程退出。
//func onServeListener(lis net.Listener) {
//	//非子进程，则什么都不走
//	if drpcGraceful.IsWorker() == false {
//		return
//	}
//	pPid := syscall.Getppid()
//	if pPid != 1 {
//		if err := graceful.SyscallKillSIGTERM(pPid); err != nil {
//			logger.Errorf("[reboot-killOldProcess] %s", err.Error())
//			return
//		}
//		logger.Infof("平滑重启中,子进程[%d]已向父进程[%d]发送信号'SIGTERM'", syscall.Getpid(), pPid)
//	}
//}

////关闭服务
//func shutdown() error {
//	endpointList.rwMu.Lock()
//
//	var list []*endpoint
//	var count int
//	var errCh = make(chan error, len(list))
//	for e := range endpointList.list {
//		list = append(list, e)
//	}
//	endpointList.rwMu.Unlock()
//
//	//异步关闭端点
//	for _, e := range list {
//		count++
//		go func(e *endpoint) {
//			errCh <- e.Close()
//		}(e)
//	}
//	var err error
//	for i := 0; i < count; i++ {
//		err = errors.Merge(err, <-errCh)
//	}
//	close(errCh)
//	return err
//}

//// SetShutdown 设置优化退出方案基本参数
//func SetShutdown(timeout time.Duration, firstFunc, beforeExitingFunc func() error) {
//	//退出之前执行的方法
//	if firstFunc == nil {
//		firstFunc = func() error {
//			return nil
//		}
//	}
//	//退出之后执行的方法
//	if beforeExitingFunc == nil {
//		beforeExitingFunc = func() error {
//			return nil
//		}
//	}
//	//设置
//	drpcGraceful.SetShutdown(
//		timeout,
//		func() error {
//			drpcGraceful.SetParentAddrList()
//			return errors.Merge(
//				firstFunc(),            //执行自定义方法
//				inherit.SetInherited(), //把监听的文件句柄数量写入环境变量，方便子进程使用
//			)
//		},
//		func() error {
//			return errors.Merge(shutdown(), beforeExitingFunc())
//			//return beforeExitingFunc()
//		})
//}

//// Shutdown 关闭服务
//func Shutdown(timeout ...time.Duration) {
//	drpcGraceful.Shutdown(timeout...)
//}
//
//// Reboot 重启服务
//func Reboot(timeout ...time.Duration) {
//	drpcGraceful.Reboot(timeout...)
//}

//const parentAddrKey = "GRACEFUL_INHERIT_LISTEN_PARENT_ADDR"
//const isChildKey = "GRACEFUL_IS_CHILD"
//
//// network:host:[host:port]
//var parentAddrList = make(map[string]map[string][]string, 2)
//var isChild = false
//var parentAddrListMutex sync.Mutex
//
////通过环境变量，初始化父进程监听的端口
////在服务启动的时候，首先从环境变量中获取父进程监听的地址端口，
////如果是首次启动，则不会获取到这些数据
////如果是优雅的无缝重启，则能通过环境变量获取到这些数据，从而复用链接，做到无缝重启
//func initParentAddrList() {
//	parentAddr := os.Getenv(parentAddrKey)
//	_ = json.Unmarshal(gconv.Bytes(parentAddr), &parentAddrList)
//	//判断当前进程是否是子进程
//	isChildEnv := os.Getenv(isChildKey)
//	if isChildEnv == "true" {
//		isChild = true
//	}
//}
//
//// 服务将要重启之前，把当前进程监听的地址端口序列化写入环境变量
//func setParentAddrList() {
//	b, _ := json.Marshal(parentAddrList)
//	env := make(map[string]string)
//	env[parentAddrKey] = gconv.String(b)
//	env[isChildKey] = "true"
//	drpcGraceful.AddInherited(nil, env)
//}
//
//// PushParentAddr 把监听的地址端口写入到变量，优雅结束的时候写入到环境变量，让子进程使用
//func PushParentAddr(network, host, addr string) {
//	parentAddrListMutex.Lock()
//	defer parentAddrListMutex.Unlock()
//	unifyLocalhost(&host)
//	nw, found := parentAddrList[network]
//	if !found {
//		nw = make(map[string][]string)
//		parentAddrList[network] = nw
//	}
//	nw[host] = append(nw[host], addr)
//}
//
//// PopParentAddr 从监听变量中出栈指定的地址端口
//func PopParentAddr(network, host, addr string) string {
//	parentAddrListMutex.Lock()
//	defer parentAddrListMutex.Unlock()
//	unifyLocalhost(&host)
//	nw, found := parentAddrList[network]
//	if !found {
//		return addr
//	}
//	h, ok := nw[host]
//	if !ok || len(h) == 0 {
//		return addr
//	}
//	nw[host] = h[1:]
//	return h[0]
//}
//
//// 针对地址格式做统一的转换
//func unifyLocalhost(host *string) {
//	switch *host {
//	case "localhost":
//		*host = "127.0.0.1"
//	case "0.0.0.0":
//		*host = "::"
//	}
//}
