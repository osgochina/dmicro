package dserver

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"github.com/osgochina/dmicro/drpc"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/supervisor/process"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

const multiProcessMasterEnv = "DServerMultiMasterProcess"

// ProcessModel 进程模型
type ProcessModel int

// ProcessModelSingle 单进程模式
const ProcessModelSingle ProcessModel = 0

// ProcessModelMulti 多进程模型
const ProcessModelMulti ProcessModel = 1

// DServer 服务对象
type DServer struct {
	name                 string // 应用名
	manager              *process.Manager
	serviceList          *gmap.TreeMap    //启动的服务列表
	started              *gtime.Time      //服务启动时间
	shutting             bool             // 服务正在关闭
	beforeStopFunc       StopFunc         //服务关闭之前执行该方法
	pidFile              string           //pid文件的路径
	sandboxNames         *garray.StrArray // 启动服务的名称
	config               *gcfg.Config     ///服务的配置信息
	inheritAddr          []InheritAddr    // 多进程模式，开启平滑重启逻辑模式下需要监听的列表
	procModel            ProcessModel     // 进程模式，processModelSingle 单进程模型，processModelMulti 多进程模型
	graceful             *graceful
	masterBool           bool //是否是主进程
	startFunction        StartFunc
	grumbleApp           *grumble.App
	cobraCmd             *cobra.Command // cobra根命令
	cobraRootCmdCallback func(rootCmd *cobra.Command)
	openCtl              bool          // 是否开启ctl功能，默认是开启，
	ctrlEndpoint         drpc.Endpoint // 作为服务提供管理接口
	ctrlSession          drpc.Session  // 作为客户端，链接到服务
}

// StartFunc 启动回调方法
type StartFunc func(svr *DServer)

// StopFunc 服务关闭回调方法
type StopFunc func(svr *DServer) bool

// newDServer  创建服务
func newDServer(name string) *DServer {
	svr := &DServer{
		name:         name,
		serviceList:  gmap.NewTreeMap(gutil.ComparatorString, true),
		sandboxNames: garray.NewStrArray(true),
		manager:      process.NewManager(),
		masterBool:   genv.GetVar(multiProcessMasterEnv, true).Bool(),
		openCtl:      true,
	}
	svr.graceful = newGraceful(svr)
	return svr
}

// SetPidFile 设置pid文件路径
func (that *DServer) SetPidFile(pidFile string) {
	that.pidFile = pidFile
}

// BeforeStop 设置服务重启方法
func (that *DServer) BeforeStop(f StopFunc) {
	that.beforeStopFunc = f
}

// ProcessModel 设置多进程模式
// 只有linux下才支持多进程模式
func (that *DServer) ProcessModel(model ProcessModel) {
	if runtime.GOOS == "linux" {
		v := gcmd.GetOptVar("model", gconv.String(model))
		that.procModel = ProcessModel(v.Int())
		return
	}
	that.procModel = ProcessModelSingle
}

// 启动服务
func (that *DServer) run(cmd *cobra.Command) {
	//判断是否是守护进程运行
	if e := that.demonize(that.config); e != nil {
		logger.Fatalf("error:%v", e)
	}
	//记录启动时间
	that.started = gtime.Now()

	// 初始化平滑重启的钩子函数
	that.initGraceful()

	if that.startFunction != nil {
		//执行业务入口函数
		that.startFunction(that)
	}
	//设置优雅退出时候需要做的工作
	that.graceful.setShutdown(15*time.Second, that.firstStop, that.beforeExiting)

	if that.isMaster() && that.openCtl {
		that.endpoint()
	}
	// 如果开启了多进程模式，并且当前进程在主进程中
	if that.procModel == ProcessModelMulti && that.isMaster() {
		that.runProcessModelMulti(cmd)
	} else {
		that.runProcessModelSingle(cmd)
	}
	// 当前进程模式是单进程模式，或者进程模式是多进程模式中的主进程，需要写入pid文件，其他进程不能写入
	if that.procModel == ProcessModelSingle || (that.procModel == ProcessModelMulti && that.isMaster()) {
		//写入pid文件
		that.putPidFile()
	}

	//答疑服务信息
	logger.Printf("%d: 服务已经初始化完成, %d 个协程被创建.", os.Getpid(), runtime.NumGoroutine())
	//监听重启信号
	that.graceful.graceSignal()
}

// 多进程模式下启动service进程
func (that *DServer) runProcessModelMulti(cmd *cobra.Command) {
	// 多进程模式下，master进程预先监听地址
	err := that.inheritListenerList()
	if err != nil {
		logger.Warningf("error:%v", err)
		os.Exit(255)
	}
	// 启动service进程
	that.serviceList.Iterator(func(_ interface{}, v interface{}) bool {
		dService := v.(*DService)
		dService.start(cmd)
		return true
	})
}

// 单进程模式下启动sandbox
func (that *DServer) runProcessModelSingle(cmd *cobra.Command) {
	// 业务进程启动sandbox
	that.serviceList.Iterator(func(_ interface{}, v interface{}) bool {
		dService := v.(*DService)
		dService.start(cmd)
		return true
	})
}

// Setup 启动服务，并执行传入的启动方法
func (that *DServer) setup(startFunction ...StartFunc) {
	if that.openCtl && len(os.Args) > 1 && os.Args[1] == "ctl" {
		// 开启ctl命令
		that.ctl()
		return
	}
	// 初始化命令行功能
	that.initCobra()
	if that.cobraRootCmdCallback != nil {
		that.cobraRootCmdCallback(that.cobraCmd)
	}
	if len(startFunction) > 0 {
		// 启动用户启动方法
		that.startFunction = startFunction[0]
	}
	_ = that.cobraCmd.Execute()
}

// AddSandBox 添加sandbox到服务中
// services 是可选，如果不传入则表示使用默认服务
func (that *DServer) AddSandBox(s ISandbox, services ...*DService) error {
	if _, found := that.searchDServiceBySandboxName(s.Name()); found {
		return gerror.Newf("Sandbox [%s] 已存在", s.Name())
	}
	var service *DService
	if len(services) > 0 {
		service = services[0]
	} else {
		s2, found := that.serviceList.Search("default")
		if !found {
			service = that.NewService("default")
		} else {
			service = s2.(*DService)
		}
	}
	err := service.addSandBox(s)
	if err != nil {
		return err
	}
	that.serviceList.Set(service.Name(), service)
	return nil
}

// 检查sandbox名字是否存在，全局sandbox名称唯一
func (that *DServer) searchDServiceBySandboxName(name string) (*DService, bool) {
	var found = false
	for _, v := range that.serviceList.Map() {
		dService := v.(*DService)
		_, found = dService.SearchSandBox(name)
		if found {
			return dService, true
		}
	}
	return nil, false
}

// Config 获取配置信息
func (that *DServer) Config() *gcfg.Config {
	return that.config
}

// StartTime 返回启动时间
func (that *DServer) StartTime() *gtime.Time {
	return that.started
}

//通过参数设置日志级别
// 日志级别通过环境默认分三个类型，开发环境，测试环境，生产环境
// 开发环境: 日志级别为 DEVELOP,标准输出打开
// 测试环境：日志级别为 INFO,除了debug日志，都会被打印，标准输出关闭
// 生产环境: 日志级别为 PRODUCT，会打印 WARN,ERRO,CRIT三个级别的日志，标准输出为关闭
// Debug开关会无视以上设置，强制把日志级别设置为ALL，并且打开标准输出。
func (that *DServer) initLogSetting(config *gcfg.Config) error {
	loggerCfg := config.GetJson("logger")
	if loggerCfg == nil {
		loggerCfg = gjson.New(nil)
	}
	env := config.GetString("ENV_NAME")
	level := loggerCfg.GetString("Level")
	logger.SetDebug(false)
	logger.SetStdoutPrint(false)
	//如果配置文件中的日志配置不存在，则判断环境变量，通过不同的环境变量，给与不同的日志级别
	if len(level) <= 0 {
		if env == "dev" || env == "develop" {
			level = "DEVELOP"
		} else if env == "test" {
			level = "INFO"
		} else {
			level = "PRODUCT"
		}
	}

	setConfig := g.Map{"level": level}

	if env == "dev" || env == "develop" {
		setConfig["stdout"] = true
		logger.SetDebug(true)
	}
	logPath := loggerCfg.GetString("Path")
	if len(logPath) > 0 {
		setConfig["path"] = logPath
	} else {
		logger.SetDebug(true)
	}

	// 如果开启debug模式，则无视其他设置
	if config.GetBool("Debug", false) {
		setConfig["level"] = "ALL"
		setConfig["stdout"] = true
		logger.SetDebug(true)
	}
	return logger.SetConfigWithMap(setConfig)
}

//守护进程
func (that *DServer) demonize(config *gcfg.Config) error {

	//判断是否需要后台运行
	daemon := config.GetBool("Daemon", false)
	if !daemon {
		return nil
	}

	if syscall.Getppid() == 1 {
		return nil
	}
	// 将命令行参数中执行文件路径转换成可用路径
	filePath := gfile.SelfPath()
	logger.Infof("Starting %s: ", filePath)
	arg0, e := exec.LookPath(filePath)
	if e != nil {
		return e
	}
	argv := make([]string, 0, len(os.Args))
	for _, arg := range os.Args {
		if arg == "--daemon" || arg == "-d" {
			continue
		}
		argv = append(argv, arg)
	}
	cmd := exec.Command(arg0, argv[1:]...)
	cmd.Env = os.Environ()
	// 将其他命令传入生成出的进程
	cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start() // 开始执行新进程，不等待新进程退出
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

//写入pid文件
func (that *DServer) putPidFile() {
	pid := os.Getpid()
	f, e := os.OpenFile(that.pidFile, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	if e != nil {
		logger.Fatalf("os.OpenFile: %v", e)
	}
	defer func() {
		_ = f.Close()
	}()
	if e := os.Truncate(that.pidFile, 0); e != nil {
		logger.Fatalf("os.Truncate: %v.", e)
	}
	if _, e := fmt.Fprintf(f, "%d", pid); e != nil {
		logger.Fatalf("Unable to write pid %d to file: %s.", pid, e)
	}
	logger.Printf("写入Pid:[%d]到文件[%s]", pid, that.pidFile)
}

// Shutdown 主动结束进程
func (that *DServer) Shutdown(timeout ...time.Duration) {
	that.graceful.shutdownSingle(timeout...)
}

// 重启或优雅的关闭服务以前，优先调用的收尾方法。收到信号后首先调用
func (that *DServer) firstStop() error {
	if that.shutting {
		return nil
	}
	that.shutting = true

	if (that.procModel == ProcessModelMulti && that.isMaster()) || that.procModel == ProcessModelSingle {
		if len(that.pidFile) > 0 && gfile.Exists(that.pidFile) {
			if e := gfile.Remove(that.pidFile); e != nil {
				logger.Errorf("os.Remove: %v", e)
			}
			logger.Printf("删除pid文件[%s]成功", that.pidFile)
		}
	}

	//结束服务前调用该方法,如果结束回调方法返回false，则中断结束
	if that.beforeStopFunc != nil && !that.beforeStopFunc(that) {
		err := gerror.New("执行完服务停止前的回调方法，该方法强制中断了服务结束流程！")
		logger.Warning(err)
		that.shutting = false
		return err
	}

	return nil
}

// 重启或优雅的关闭服务以前，处理完基本操作后，做最后的收尾工作，收尾完成就结束进程
func (that *DServer) beforeExiting() error {
	if that.procModel == ProcessModelMulti && that.isMaster() {
		return nil
	}
	//结束各组件
	that.serviceList.Iterator(func(_ interface{}, v interface{}) bool {
		dService := v.(*DService)
		dService.stop()
		return true
	})
	return nil
}

// IsMaster 判断当前进程是否是主进程
func (that *DServer) isMaster() bool {
	return that.masterBool
}

// NewService 创建新的服务
// 注意: 如果是多进程模式，则每个service表示一个进程
func (that *DServer) NewService(name string) *DService {
	return newDService(name, that)
}

// SetInheritListener 多进程模式下，设置要被继承的监听地址
func (that *DServer) SetInheritListener(address []InheritAddr) {
	if that.isMaster() {
		that.inheritAddr = address
	}
}
