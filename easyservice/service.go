package easyservice

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/gracefulv2"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

// EasyService 服务对象
type EasyService struct {
	sList          *gmap.IntAnyMap //启动的服务列表
	started        *gtime.Time     //服务启动时间
	shutting       bool            // 服务正在关闭
	beforeStopFunc StopFunc        //服务关闭之前执行该方法
	pidFile        string          //pid文件的路径
	processName    string          // 进程名字
	cmdParser      *gcmd.Parser    //命令行参数解析信息
	config         *gcfg.Config    ///服务的配置信息
}

// StartFunc 启动回调方法
type StartFunc func(service *EasyService)

// StopFunc 服务关闭回调方法
type StopFunc func(service *EasyService) bool

// NewEasyService  创建服务
func NewEasyService(processName ...string) *EasyService {
	svr := &EasyService{
		sList: gmap.NewIntAnyMap(true),
	}
	if len(processName) > 0 {
		svr.processName = processName[0]
	}
	return svr
}

// SetPidFile 设置pid文件
func (that *EasyService) SetPidFile(pidFile string) {
	that.pidFile = pidFile
}

// SetProcessName 设置进程名字
func (that *EasyService) setProcessName(processName string) {
	that.processName = processName
}

// BeforeStop 设置服务重启方法
func (that *EasyService) BeforeStop(f StopFunc) {
	that.beforeStopFunc = f
}

func (that *EasyService) SetGracefulModel(model gracefulv2.GracefulModel, address ...[]string) {
	graceful := gracefulv2.GetGraceful()
	graceful.SetModel(model)
	if model == gracefulv2.GracefulMasterWorker {
		//graceful.InheritedListener()
	}
}

// Setup 启动服务，并执行传入的启动方法
func (that *EasyService) Setup(startFunction StartFunc) {
	//解析命令行
	parser, err := gcmd.Parse(defaultOptions)
	if err != nil {
		logger.Fatal(err)
	}
	//解析参数
	if !that.parserArgs(parser) {
		return
	}
	//解析配置文件
	that.parserConfig(parser)

	//启动时间
	that.started = gtime.Now()

	that.cmdParser = parser

	if that.config != nil {
		//判断是否是守护进程运行
		if e := that.demonize(that.config); e != nil {
			logger.Fatalf("error:%v", e)
		}
		//初始化日志配置
		if e := that.initLogSetting(that.config); e != nil {
			logger.Fatalf("error:%v", e)
		}
	}

	//启动自定义方法
	startFunction(that)

	//写入pid文件
	that.putPidFile()

	that.sList.Iterator(func(k int, v interface{}) bool {
		sandbox := v.(ISandBox)
		go func() {
			e := sandbox.Setup()
			if e != nil {
				logger.Warning(e)
			}
		}()
		return true
	})

	//设置优雅退出时候需要做的工作
	gracefulv2.GetGraceful().SetShutdown(15*time.Second, that.firstSweep, that.beforeExiting)
	//等待服务结束
	logger.Noticef("服务已经初始化完成, %d 个协程被创建.", runtime.NumGoroutine())
	//设置进程名
	if len(that.processName) > 0 {
		setProcessName(that.processName)
	}
	////发送启动成功信号
	gracefulv2.GetGraceful().OnStart()
	//监听重启信号
	gracefulv2.GetGraceful().GraceSignal()
}

func (that *EasyService) AddSandBox(s ISandBox) {
	that.sList.Set(s.ID(), s)
}

// GetSandBox 获取指定的服务沙盒
func (that *EasyService) GetSandBox(id int) ISandBox {
	s, found := that.sList.Search(id)
	if !found {
		return nil
	}
	return s.(ISandBox)
}

// Config 获取配置信息
func (that *EasyService) Config() *gcfg.Config {
	return that.config
}

// CmdParser 返回命令行解析
func (that *EasyService) CmdParser() *gcmd.Parser {
	return that.cmdParser
}

// StartTime 返回启动时间
func (that *EasyService) StartTime() *gtime.Time {
	return that.started
}

//设置日志级别
func (that *EasyService) initLogSetting(config *gcfg.Config) error {
	level := config.GetString("logger.Level", "PRODUCT")

	env := that.config.GetString("ENV_NAME")
	if len(env) > 0 && env == "dev" || env == "develop" {
		level = "DEVELOP"
	}
	err := logger.SetConfigWithMap(g.Map{
		"path":   config.GetString("logger.Path"),
		"level":  level,
		"stdout": true,
	})
	if err != nil {
		return err
	}
	// 开启debug模式
	debug := config.GetBool("Debug", false)
	if debug {
		_ = logger.SetLevelStr("ALL")
	}
	return nil
}

//守护进程
func (that *EasyService) demonize(config *gcfg.Config) error {

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
func (that *EasyService) putPidFile() {
	f, e := os.OpenFile(that.pidFile, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	if e != nil {
		logger.Fatalf("os.OpenFile: %v", e)
	}
	if e := os.Truncate(that.pidFile, 0); e != nil {
		logger.Fatalf("os.Truncate: %v.", e)
	}
	if _, e := fmt.Fprintf(f, "%d", os.Getpid()); e != nil {
		logger.Fatalf("Unable to write pid %d to file: %s.", os.Getpid(), e)
	}
}

// Shutdown 主动结束进程
func (that *EasyService) Shutdown(timeout ...time.Duration) {
	//drpc.Shutdown(timeout...)
}

func (that *EasyService) firstSweep() error {
	if that.shutting {
		return nil
	}
	that.shutting = true
	//结束服务前调用该方法,如果结束回调方法返回false，则中断结束
	if that.beforeStopFunc != nil && !that.beforeStopFunc(that) {
		err := gerror.New("执行完服务停止前的回调方法，该方法强制中断了服务结束流程！")
		logger.Info(err)
		that.shutting = false
		return err
	}
	if len(that.pidFile) > 0 && gfile.Exists(that.pidFile) {
		if e := gfile.Remove(that.pidFile); e != nil {
			logger.Errorf("os.Remove: %v", e)
		}
		logger.Infof("删除pid文件[%s]成功", that.pidFile)
	}
	return nil
}

//进行结束收尾工作
func (that *EasyService) beforeExiting() error {
	//结束各组件
	that.sList.Iterator(func(k int, v interface{}) bool {
		service := v.(ISandBox)
		if e := service.Shutdown(); e != nil {
			logger.Errorf("Service %s .Shutdown: %v", service.Name(), e)
		} else {
			logger.Infof("%s Service Stoped.", service.Name())
		}
		return true
	})
	return nil
}
