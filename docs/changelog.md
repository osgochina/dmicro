# 版本更新记录

## v1.1.0 (2022-09-07)
1. 替换 `dserver` 组件的命令行解析工具为 `Cobra`组件。
2. `dserver` 组件中的 `ctrl` 命令替换为 `ctl`
3. `dserver` 组件中的 `rpc server option`参数变成 `ServerName`
4. 取消`endpoint`与`eventbus`强关联。
5. 新增 `Metrics` 指标组件,支持统计运行指标。
6. 修复 `rpc client` 与 `rpc server` 中 `close` 报错的问题.
7. 移除 `EndpiontConfig`中的`CountTime`参数，默认开启请求耗时统计。
8. 完善 `Rpc Server`,`Rpc Client`,`Metrics`组件的使用文档。

## v1.0.1 (2022-08-13)
1. 修复dServer在macos下报错的问题。
2. 增加`memory registry` 服务注册组件。
3. 优化`RPC Client` 接口。
4. 优化`RPC Server` 接口。
5. 完善文档，增加更多的使用示例。

## v1.0.0 (2022-08-01)
1. 发布新组件`DServer`,该组件是`easyserver`组件的升级版本。
   1. `dserver`服务管理功能能够让你专注于编写业务代码，编译部署后的运行时管理就交给它吧。
   2. 支持单进程，多进程模式，单进程模式方便开发，多进程模式适合业务隔离。
   3. 原生支持平滑重启功能。
   4. 方便的扩展命令行功能。
   5. 原生支持命令行ctl，方便开启关闭服务，重启服务，开启debug模式，查看实时运行日志，查看运行指标。
   6. 更多`DServer`组件的介绍，请看[文档](./dserver/readme.md)
2. `supervisor`组件`api`大改,从功能独立的组件融合进框架，更好的与`dServer`组合。
3. `drpc`组件修复`unix socket`链接的监听。
4. 增加`benchmark`测试用例。
5. 完善文档，增加更多的使用示例。


## v0.6.2 (2022-06-14)
1. 完成服务注册功能`registry`
2. 完成服务发下功能`selector`
3. 增加`mdns`服务发现组件
4. 升级gf依赖版本至`v1.16.9`
5. 改造`rpc client`支持服务发现功能。
6. `easy service`暴露出`Help`，`Version`方法，方便业务调用。

## v0.5.0 (2022-01-19)
1. 增加`build.sh`编译脚本，支持设置编译变量，方便使用`easyservice`组件使用`version`命令展示编译信息。
2. 完善`easyservice`的日志级别配置.
3. 移除`easyservice`的默认`network`,`host`,`port`参数支持。
4. 增加`master-worker`进程模型下，worker进程异常退出后，master进程自动拉起功能。
5. 修复了`reply`消息解包失败不能正确报错的问题。
6. 修复了`jsonproto`协议不能正确处理字符串的问题。
7. 增加drpc的内部日志组件，支持重设日志组件，方便与默认的日志组件区分。

## v0.4.0 (2022-01-10)
1. 优化`easyService`服务的行为,增加`-c,--config`参数的支持.
2. 支持`easyservice`的quit命令.
3. 修复`window`，`macos`系统下不支持`quic`和`kcp`的问题.
4. 新增兼容多平台的`signal`发送组件.
5. 支持`Websocket`协议.
6. 增加并发请求客户端`Multiclient`.

## v0.3.0 (2021-12-31)
1. 修复在dial成功后，对端马上关闭了链接，造成重试死循环的bug，增加最大重试次数.
2. 修复平滑重启中监听地址0.0.0.0不起效的问题.
3. 增加`proxy`插件的的文档及测试用例.
4. `easyservice`组件进程退出之前先删除pid文件。
5. `easyservice`支持从配置文件中读取sandbox的id.
6. 支持`quic`协议.
7. 支持`kcp`协议.

## v0.2.0 (2021-12-20)
1. 支持`ProtoBuf`协议.
2. 增加安全传输`SecureBodyPlugin`插件.
3. 完善文档，增加框架`logo`.

## v0.1.4 (2021-12-06)
1. MacOS 支持
2. Windows 支持
3. `平滑重启` 逻辑支持，支持`父子进程`模式以及`Master-Worker`模式。

## v0.0.8 (2021-10-22)

1. 完善`平滑重启`的逻辑。
2. `Supervisor进程监控管理`模块新增接口，支持从`supervisor.ini`格式的配置文件中载入配置，从而启动配置的进程。

## v0.0.7

1. 完成`EventBus(事件总线)`功能的开发，详情请见 [eventBus(事件总线)](component/eventBus.md)
2. 修复`Supervisor进程监控管理`模块的一些bug。

## v0.0.6

1. 增加`Supervisor进程监控管理`模块，初步支持管理多个子进程，后续会支持通过`supervisor.conf`配置文件配置子进程启动。
2. 完善`EasyService`快速服务创建模块。
3. 完善`drpc message`模块的测试用例及文档。

## v0.0.5

1. 增加 `BeforeCloseEndpoint` 关闭Endpoint之前触发该事件.
2. 增加 `AfterCloseEndpointPlugin` 关闭Endpoint之后触发该事件.
3. 完善 `event`事件的测试用例。
4. 完善 `confg`配置的测试用例。