# 插件开发

`drpc`的插件开发非常方便，只需要熟悉了`Hook`点的触发时间，就能根据需求，开发满足需要的插件。

以下代码展示了如何开发一个简单的插件。

```go

type myPlugin struct {
	newParams string
	myParams  string
}

// 强制检查`myPlugin`是否实现了对应的hook节点，为了防止hook点太多的时候遗漏
var (
    _ drpc.AfterNewEndpointPlugin    = new(myPlugin)
    _ drpc.BeforeCloseEndpointPlugin = new(myPlugin)
    _ drpc.AfterDialPlugin           = new(myPlugin)
    _ drpc.AfterAcceptPlugin         = new(myPlugin)
)


func NewMyPlugin(newParams string) *myPlugin {
	return &myPlugin{newParams: newParams}
}

func (that *myPlugin) Name() string {
	return "myPlugin"
}

// AfterNewEndpoint endpoint创建成功后调用
func (that *myPlugin) AfterNewEndpoint(drpc.EarlyEndpoint) error {
	that.myParams = "myParams"
	fmt.Println("AfterNewEndpoint")
	return nil
}

// BeforeCloseEndpoint endpoint关闭之前调用
func (that *myPlugin) BeforeCloseEndpoint(drpc.Endpoint) error {
	fmt.Println("BeforeCloseEndpoint")
	return nil
}

// AfterDial 客户端链接到服务端成功后调用(endpoint作为客户端时生学校)
func (that *myPlugin) AfterDial(sess drpc.EarlySession, isRedial bool) *drpc.Status {
	fmt.Println("AfterDial")
	fmt.Println(that.myParams)
	fmt.Println(sess.RemoteAddr())
	return nil
}

// AfterAccept 服务端接收客户端请求成功后调用(endpoint作为服务端时生效)
func (that *myPlugin) AfterAccept(sess drpc.EarlySession) *drpc.Status {
	fmt.Println("AfterAccept")
	fmt.Println(that.myParams)
	fmt.Println(sess.RemoteAddr())
	return nil
}
```

我们从头开始讲解开发`myPlugin`插件的流程。

1. 定义插件结构体
2. 创建结构体New函数，当然，也可以直接`&myPlugin`。
3. 实现`Name()`方法，定义插件的名字，插件名字必须唯一。
4. 实现`Hook`点方法，运行到对应的阶段，会调用该方法。
5. 把`Plugin`对象传入`Endpoint`.

以上代码可以放入`main.go`运行。

```go

func main() {
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	},
		NewMyPlugin("clownfish"),
	)
	go svr.ListenAndServe()

	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{PrintDetail: false},
		NewMyPlugin("clownfish"),
	)
	cli.Dial(":9090")
	time.Sleep(time.Second * 20)
}
```

编译运行
```shell
$ go build main.go
$ ./main

AfterNewEndpoint
2022-07-26 15:23:39.294 pid:772747,启动监听并提供服务：(network:tcp, addr:[::]:9090) 
AfterNewEndpoint
AfterDial
myParams
127.0.0.1:9090
2022-07-26 15:23:40.295 [NOTI] dial ok (network:tcp, addr::9090, id:127.0.0.1:47089) 
AfterAccept
myParams
127.0.0.1:47089
2022-07-26 15:23:40.295 [NOTI] accept ok (network:tcp, addr:127.0.0.1:47089, id:127.0.0.1:47089) 
```