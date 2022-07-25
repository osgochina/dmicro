# 插件系统 plugin

插件系统给框架带来了极大的扩展性和灵活性，是整个框架的一个灵魂模块，有了它，框架就有了无限可能。

什么样的插件系统才能算是优雅呢？我能想到的有以下几点：

- 合理且丰富的`hook`位置,能够覆盖整个框架的生命周期，贯穿通讯的各个环节。
- 每个`hook`位置的入参和出参都是经过精心设计。
- 每个插件都能够使用多个`hook`位置，每个`hook`位置都能被多个插件使用。
- 设计的足够简洁，优雅。能方便的进行二次开发定制。

### 插件使用

插件使用很简单,可以在创建`endpoin`的时候为它载入`Plugin`对象,也可以通过`endpoint.PluginContainer()`获取插件容器，添加删除插件。

#### 一、创建`endpoint`对象的时候加载插件

` drpc.NewEndpoint()`方法的第二个参数是可变参数，后面传入的都是插件对象，可以是框架内置，也可以是自定义插件，可以传入一个或多个插件对象。
```go
// 心跳插件的使用
func main() {
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	},
		heartbeat.NewPong(), // 使用pong插件，响应客户端的ping请求
	)

	go svr.ListenAndServe()

	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{PrintDetail: true},
		heartbeat.NewPing(3, true), // 使用ping插件，发送ping请求给服务端
	)
	cli.Dial(":9090")
	time.Sleep(time.Second * 20)
}
```

#### 二、通过`endpoint.PluginContainer()`管理插件

`endpoint.PluginContainer()` 有5个方法

| 方法名 | 作用                        |
|-----|---------------------------|
|AppendLeft| 往插件容器的左边添加插件(执行插件的时候优先执行) |
|AppendRight| 往插件容器的右边添加插件(执行插件的时候最后执行) |
|GetAll| 获取所有插件对象                  |
|GetByName| 根据插件名字获取插件对象              |
|Remove| 移除指定名称的插件                 |

因为每个`hook`点能够注册多个插件，那么插件就必定有执行的先后顺序。

- 具体的顺序在创建`endpoint`对象的时候是按参数的先后顺序。
- 通过`PluginContainer()`对象添加的时候就需要调用不同的方法来决定。

```go
func main() {
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	},
	)
	svr.PluginContainer().AppendRight(ignorecase.NewIgnoreCase()) //使用IgnoreCase插件，活跃serviceName的大小写
	go svr.ListenAndServe()

	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{PrintDetail: false},
	)
	cli.PluginContainer().AppendRight(ignorecase.NewIgnoreCase()) // 使用IgnoreCase插件，活跃serviceName的大小写
	cli.Dial(":9090")
	time.Sleep(time.Second * 20)
}
```


#### 三、通过`endpoint.PluginContainer()`添加插件的奇怪现象
```go
func main() {
	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		ListenPort:  9090,
		PrintDetail: true,
	},
	)
	svr.PluginContainer().AppendRight(heartbeat.NewPong()) //使用pong插件，响应客户端的ping请求
	go svr.ListenAndServe()

	time.Sleep(time.Second)

	cli := drpc.NewEndpoint(
		drpc.EndpointConfig{PrintDetail: true},
	)
	cli.PluginContainer().AppendRight(heartbeat.NewPing(3, true)) // 使用ping插件，发送ping请求给服务端
	cli.Dial(":9090")
	time.Sleep(time.Second * 20)
}
```
如果你使用以上代码运行，会发现并没有如预期那般发送心跳包，难道是这种方法添加插件不好用吗？

并不是，要理解这个问题的原因，需要理解插件使用的`hook`点，`heartbeat.Ping`插件启动使用到的`hook`点是`AfterNewEndpoint`,根据`hook`[文档](drpc/hook?id=endpoint生命周期中会触发的钩子列表)所展示的内容可知，
触发时间点是`创建Endpoint之后触发`,因为`NewEndpoint`的时候已经触发该`hook`点，所以再通过`PluginContainer`添加插件是不会其效果的。

### 框架内置的插件

| 插件名字       | 插件作用                |
|------------|---------------------|
| IgnoreCase | 忽略`serviceName`的大小写 |
| Heartbeat  | 心跳插件，保持两个服务之间的链接    |
| Auth       | 票据检查，服务之间连接权限校验     |
| SecureBody | 消息内容加密              |
| Proxy      | 请求代理                |
| Event      | 事件总线                |

