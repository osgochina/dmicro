package dserver

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"reflect"
	"unsafe"
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

func (that *DService) addSandBox(s ISandbox) error {
	name := s.Name()
	_, found := that.sList.Search(name)
	if found {
		return gerror.Newf("Sandbox [%s] 已存在", name)
	}
	p, err := that.makeSandBox(s)
	if err != nil {
		return err
	}
	//赋值
	p.handlerCtx.service = that
	that.sList.Set(p.Name(), p)
	return nil
}

// Iterator 迭代sandbox
func (that *DService) iterator(f func(name string, sandbox ISandbox) bool) {
	that.sList.Iterator(func(k string, v interface{}) bool {
		return f(k, v.(ISandbox))
	})
}

func (that *DService) Name() string {
	return that.name
}

// SearchSandbox 获取指定的sandbox
func (that *DService) SearchSandbox(name string) (ISandbox, bool) {
	val, found := that.sList.Search(name)
	if found {
		return val.(ISandbox), true
	}
	return nil, false
}

// 通过反射生成私有sandbox对象
func (that *DService) makeSandBox(s ISandbox) (*privateSandbox, error) {
	var (
		cType = reflect.TypeOf(s)
	)
	//判断是否是指针类型
	if cType.Kind() != reflect.Ptr {
		return nil, gerror.Newf("make sandbox: the type is not struct point: %s", cType.String())
	}
	var cTypeElem = cType.Elem()
	//判断是否是struct类型
	if cTypeElem.Kind() != reflect.Struct {
		return nil, gerror.Newf("make sandbox: the type is not struct point: %s", cType.String())
	}
	//如果结构体没有实现 ISandbox 的方法，或者不是匿名结构体
	iType, ok := cTypeElem.FieldByName("SandboxCtx")
	if !ok || !iType.Anonymous {
		return nil, gerror.Newf("make sandbox: the struct do not have anonymous field dserver.SandboxCtx: %s", cType.String())
	}
	var callCtxOffset = iType.Offset

	ctrl := reflect.New(cTypeElem)
	privateBox := &privateSandbox{
		original: s,
		sandboxVal: &sandboxReflectValue{
			ctrl: ctrl,
			//这种写法参考https://blog.csdn.net/u010853261/article/details/103826830中的模式三
			//将非类型安全指针转换为一个uintptr值，然后此uintptr值参与各种算术运算，再将算术运算的结果uintptr值转回非类型安全指针
			ctxPtr: (*SandboxCtx)(unsafe.Pointer(uintptr(unsafe.Pointer(ctrl.Pointer())) + callCtxOffset)),
		},
		cType:      cType,
		handlerCtx: new(handlerCtx),
	}

	return privateBox, nil
}
