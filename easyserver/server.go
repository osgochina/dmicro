package easyserver

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/osgochina/dmicro/logger"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

// Server 服务对象
type Server struct {
	sList          *gmap.IntAnyMap //启动的服务列表
	started        *gtime.Time     //服务启动时间
	shutting       bool            // 服务正在关闭
	bus            chan struct{}   // 控制总线
	beforeStopFunc StopFunc        //服务关闭之前执行该方法
	pidFile        string          //pid文件的路径
	processName    string          // 进程名字
	cmdParser      *gcmd.Parser    //命令行参数解析信息
	config         *gcfg.Config    ///服务的配置信息
}

// StartFunc 启动回调方法
type StartFunc func(service *Server)

// StopFunc 服务关闭回调方法
type StopFunc func(service *Server) bool

// NewServer 创建服务
func NewServer(processName ...string) *Server {
	svr := &Server{
		processName: "default-server",
		sList:       gmap.NewIntAnyMap(true),
	}
	if len(processName) > 0 {
		svr.processName = processName[0]
	}
	return svr
}

// SetPidFile 设置pid文件
func (that *Server) SetPidFile(pidFile string) {
	that.pidFile = pidFile
}

// SetProcessName 设置进程名字
func (that *Server) SetProcessName(processName string) {
	that.processName = processName
}

// BeforeStop 设置服务重启方法
func (that *Server) BeforeStop(f StopFunc) {
	that.beforeStopFunc = f
}

// Setup 启动服务，并执行传入的启动方法
func (that *Server) Setup(startFunction StartFunc) {
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
		if e := sandbox.Setup(); e != nil {
			logger.Fatalf("Service.Int: %v.", e)
		}
		return true
	})
	//监听重启信号
	go that.signaler()
	//等待服务结束
	logger.Noticef("服务已经初始化完成, %d 个协程被创建.\n", runtime.NumGoroutine())
	that.master()
}

func (that *Server) AddSandBox(s ISandBox) {
	that.sList.Set(s.ID(), s)
}

// GetSandBox 获取指定的服务沙盒
func (that *Server) GetSandBox(id int) ISandBox {
	s, found := that.sList.Search(id)
	if !found {
		return nil
	}
	return s.(ISandBox)
}

// Config 获取配置信息
func (that *Server) Config() *gcfg.Config {
	return that.config
}

// CmdParser 返回命令行解析
func (that *Server) CmdParser() *gcmd.Parser {
	return that.cmdParser
}

// StartTime 返回启动时间
func (that *Server) StartTime() *gtime.Time {
	return that.started
}

//设置日志级别
func (that *Server) initLogSetting(config *gcfg.Config) error {
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
		logger.SetDebug(debug)
	}
	return nil
}

//守护进程
func (that *Server) demonize(config *gcfg.Config) error {

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
func (that *Server) putPidFile() {
	f, e := os.OpenFile(that.pidFile, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	if e != nil {
		logger.Fatalf("os.OpenFile: %v", e)
	}
	if e := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		logger.Fatalf("syscall.Flock: %v, already in running?", e)
	}
	if e := os.Truncate(that.pidFile, 0); e != nil {
		logger.Fatalf("os.Truncate: %v.", e)
	}
	if _, e := fmt.Fprintf(f, "%d", os.Getpid()); e != nil {
		logger.Fatalf("Unable to write pid %d to file: %s.", os.Getpid(), e)
	}
}

//监听信号
func (that *Server) signaler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGPROF) //增加SIGPROF信号支持
	for {
		switch <-ch {
		case syscall.SIGHUP:
			logger.Notice("Received SIGHUP signal, configuring service...")
		case syscall.SIGTERM, syscall.SIGINT:
			logger.Notice("Received SIGTERM signal, shutting down...")
			that.Shutdown()
		}
	}
}

// 主协程阻塞等待结束的信号
func (that *Server) master() {
	that.bus = make(chan struct{})
	for {
		select {
		case s, ok := <-that.bus:
			if !ok {
				logger.Noticef("%s Service exit properly.", that.processName)
				return
			}
			logger.Noticef("%s received signal: %v", that.processName, s)
			return
		}
	}
}

// Shutdown 关闭服务
func (that *Server) Shutdown() {
	if that.shutting {
		return
	}
	that.shutting = true
	//结束服务前调用该方法,如果结束回调方法返回false，则中断结束
	if that.beforeStopFunc != nil && !that.beforeStopFunc(that) {
		logger.Debug("执行完服务停止前的回调方法，该方法强制中断了服务结束流程！")
		that.shutting = false
		return
	}
	go that.shutdownNext()
}

//进行结束收尾工作
func (that *Server) shutdownNext() {

	//初始化各个服务组件
	that.sList.Iterator(func(k int, v interface{}) bool {
		service := v.(ISandBox)
		if e := service.Shutdown(); e != nil {
			logger.Errorf("Service %s .Shutdown: %v", service.Name(), e)
		} else {
			logger.Infof("%s Service Stoped.", service.Name())
		}
		return true
	})
	// remove the pidFile
	if len(that.pidFile) > 0 {
		if e := gfile.Remove(that.pidFile); e != nil {
			logger.Errorf("os.Remove: %v\n", e)
		}
		logger.Infof("Remove pidFile %s successful\n", that.pidFile)
	}
	//结束主协程
	that.bus <- struct{}{}
}
