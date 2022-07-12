package dserver

import (
	"context"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"reflect"
)

type DService struct {
	server *DServer
	name   string
	sList  *gmap.StrAnyMap //启动的服务列表
}

func newDService(name string, server *DServer) *DService {
	return &DService{
		name:   name,
		server: server,
		sList:  gmap.NewStrAnyMap(true),
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
		return s.(ISandbox), true
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
	that.sList.Set(s1.Name(), s1)
	return nil
}

// 移除sandbox
func (that *DService) removeSandbox(name string) {
	that.sList.Remove(name)
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
