package dserver

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/supervisor/process"
	"os"
	"reflect"
	"time"
)

// DService 服务对象，每个DService中可以存在多个sandbox
// 每个DServer中可以存在多个DService
type DService struct {
	server *DServer
	name   string
	sList  *gmap.TreeMap //启动的服务列表
}

// 创建DService对象
func newDService(name string, server *DServer) *DService {
	return &DService{
		name:   name,
		server: server,
		sList:  gmap.NewTreeMap(gutil.ComparatorString, true),
	}
}

// Name 获取服务名
func (that *DService) Name() string {
	return that.name
}

// SearchSandBox 搜索同一个服务下的其他sandbox
func (that *DService) SearchSandBox(name string) (ISandbox, bool) {
	s, found := that.sList.Search(name)
	if found {
		return s.(*sandboxContainer).sandbox, true
	}
	return nil, false
}

func (that *DService) addSandBox(s ISandbox) error {
	name := s.Name()
	_, found := that.sList.Search(name)
	if found {
		return gerror.Newf("Sandbox [%s] 已存在", name)
	}
	s1, err := that.makeSandBox(s)
	if err != nil {
		return err
	}
	that.sList.Set(s1.Name(), &sandboxContainer{
		sandbox: s1,
		state:   process.Unknown,
	})
	return nil
}

// 启动该service
func (that *DService) start(c *grumble.Context) {
	if that.server.procModel == ProcessModelMulti && that.server.isMaster() {
		if that.sList.Size() == 0 {
			return
		}
		// 如果命令行传入了需要启动的服务名称，则需要把改服务名提取出来，作为启动参数
		var sandBoxNames []string
		if that.server.sandboxNames.Len() > 0 {
			for _, name := range that.sList.Keys() {
				if that.server.sandboxNames.ContainsI(gconv.String(name)) {
					sandBoxNames = append(sandBoxNames, gconv.String(name))
				}
			}
		} else {
			sandBoxNames = gconv.Strings(that.sList.Keys())
		}
		// 如果未匹配服务名称，则说明该service不需要启动
		if len(sandBoxNames) == 0 {
			return
		}
		var args = []string{"start"}

		if len(that.server.config.GetString("ENV_NAME")) > 0 {
			args = append(args, fmt.Sprintf("--env=%s", that.server.config.GetString("ENV_NAME")))
		}
		confFile := c.Flags.String("config")
		if len(confFile) > 0 {
			args = append(args, fmt.Sprintf("--config=%s", confFile))
		}
		if that.server.config.GetBool("Debug") {
			args = append(args, "--debug")
		}
		args = append(args, sandBoxNames...)
		p, e := that.server.manager.NewProcessByOptions(process.NewProcOptions(
			process.ProcCommand(os.Args[0]),
			process.ProcName(that.Name()),
			process.ProcArgs(args...),
			process.ProcSetEnvironment(isChildKey, "true"),
			process.ProcSetEnvironment(multiProcessMasterEnv, "false"),
			process.ProcStdoutLog("/dev/stdout", ""),
			process.ProcRedirectStderr(true),
			process.ProcAutoReStart(process.AutoReStartTrue),             // 自动重启
			process.ProcExtraFiles(that.server.graceful.getExtraFiles()), // 与获取inheritedEnv的顺序不能错乱
			process.ProcEnvironment(that.server.graceful.inheritedEnv.Map()),
			process.ProcStopSignal("SIGQUIT", "SIGTERM"), // 退出信号
			process.ProcStopWaitSecs(int(minShutdownTimeout/time.Second)),
		))
		if e != nil {
			logger.Warning(e)
		}
		p.Start(true)
		return
	}

	for name, sandbox := range that.sList.Map() {
		s := sandbox.(*sandboxContainer)
		// 如果命令行传入了要启动的服务名，则需要匹配启动对应的sandbox
		if that.server.sandboxNames.Len() > 0 && !that.server.sandboxNames.ContainsI(s.sandbox.Name()) {
			that.removeSandbox(gconv.String(name))
			return
		}
		s.started = gtime.Now()
		s.state = process.Running
		go func(s1 *sandboxContainer) {
			e := s1.sandbox.Setup()
			if e != nil && s1.state != process.Stopping {
				s1.state = process.Stopped
				logger.Warningf("Sandbox Setup Return: %v", e)
			}
		}(s)
	}
}

// 关闭该server
func (that *DService) stop() {
	for _, sandbox := range that.sList.Map() {
		s := sandbox.(*sandboxContainer)
		if s.state == process.Running {
			s.state = process.Stopping
			if e := s.sandbox.Shutdown(); e != nil {
				logger.Errorf("服务 %s .结束出错，error: %v", s.sandbox.Name(), e)
			} else {
				logger.Printf("%s 服务 已结束.", s.sandbox.Name())
			}
			s.state = process.Stopped
			s.stopTime = gtime.Now()
		}
	}
	return
}

// 启动指定的sandbox
func (that *DService) startSandbox(name string) error {
	s, found := that.sList.Search(name)
	if !found {
		return fmt.Errorf("未找到[%s]", name)
	}
	sc := s.(*sandboxContainer)
	if sc.state == process.Starting || sc.state == process.Running {
		return fmt.Errorf("sandbox[%s]正在运行中", name)
	}
	sc.started = gtime.Now()
	sc.state = process.Running
	go func(s1 *sandboxContainer) {
		e := s1.sandbox.Setup()
		if e != nil && s1.state != process.Stopping {
			s1.state = process.Stopped
			logger.Warningf("Sandbox Setup Return: %v", e)
		}
	}(sc)
	return nil
}

// 关闭指定的sandbox
func (that *DService) stopSandbox(name string) error {
	s, found := that.sList.Search(name)
	if !found {
		return fmt.Errorf("未找到[%s]", name)
	}
	sc := s.(*sandboxContainer)
	if sc.state == process.Running {
		sc.state = process.Stopping
		err := sc.sandbox.Shutdown()
		sc.state = process.Stopped
		sc.stopTime = gtime.Now()
		return err
	}
	return nil
}

// 移除sandbox
func (that *DService) removeSandbox(name string) {
	value := that.sList.Remove(name)
	if value == nil {
		return
	}
	sandbox := value.(*sandboxContainer)
	if sandbox.state == process.Running {
		err := that.stopSandbox(name)
		if err != nil {
			logger.Error(err)
		}
	}
}

// 通过反射生成私有sandbox对象
func (that *DService) makeSandBox(s ISandbox) (ISandbox, error) {
	var (
		cType  = reflect.TypeOf(s)
		cValue = reflect.ValueOf(s)
	)
	//判断是否是指针类型
	if cType.Kind() != reflect.Ptr {
		return nil, gerror.Newf("生成Sandbox: 传入的Sandbox对象不是指针类型: %s", cType.String())
	}
	var cTypeElem = cType.Elem()
	//判断是否是struct类型
	if cTypeElem.Kind() != reflect.Struct {
		return nil, gerror.Newf("生成Sandbox: 传入的Sandbox对象不是struct类型: %s", cType.String())
	}
	//如果结构体没有实现 SandboxCtx 的方法，或者不是匿名结构体
	iType, ok := cTypeElem.FieldByName("BaseSandbox")
	if !ok || !iType.Anonymous {
		return nil, gerror.Newf("生成Sandbox: 传入的Sandbox对象未继承 dserver.BaseSandbox : %s", cType.String())
	}

	_, found := cType.MethodByName("Setup")
	if !found {
		return nil, gerror.Newf("生成Sandbox: 传入的Sandbox对象未实现Setup方法")
	}

	_, found = cType.MethodByName("Shutdown")
	if !found {
		return nil, gerror.Newf("生成Sandbox: 传入的Sandbox对象未实现Shutdown方法")
	}

	_, found = cType.MethodByName("Name")
	if !found {
		return nil, gerror.Newf("生成Sandbox: 传入的Sandbox对象未实现Name方法")
	}
	iValue := cValue.Elem().FieldByName("Service")
	if iValue.CanSet() {
		iValue.Set(reflect.ValueOf(that))
	}
	iValue = cValue.Elem().FieldByName("Context")
	if iValue.CanSet() {
		iValue.Set(reflect.ValueOf(context.Background()))
	}
	iValue = cValue.Elem().FieldByName("Config")
	if iValue.CanSet() {
		c := &Config{}
		c.Config = that.server.config
		iValue.Set(reflect.ValueOf(c))
	}
	return s, nil
}
