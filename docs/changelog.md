## 更新日志

### v0.4.0 (2022-01-10)
1. 优化`easyService`服务的行为,增加`-c,--config`参数的支持.
2. 支持`easyservice`的quit命令.
3. 修复`window`，`macos`系统下不支持`quic`和`kcp`的问题.
4. 新增兼容多平台的`signal`发送组件.
5. 支持`Websocket`协议.
6. 增加并发请求客户端`Multiclient`.

### v0.3.0 (2021-12-31)
1. 修复在dial成功后，对端马上关闭了链接，造成重试死循环的bug，增加最大重试次数.
2. 修复平滑重启中监听地址0.0.0.0不起效的问题.
3. 增加`proxy`插件的的文档及测试用例.
4. `easyservice`组件进程退出之前先删除pid文件。
5. `easyservice`支持从配置文件中读取sandbox的id.
6. 支持`quic`协议.
7. 支持`kcp`协议.

### v0.2.0 (2021-12-20)
1. 支持`ProtoBuf`协议.
2. 增加安全传输`SecureBodyPlugin`插件.
3. 完善文档，增加框架`logo`.

### v0.1.4 (2021-12-06)
1. MacOS 支持
2. Windows 支持
3. `平滑重启` 逻辑支持，支持`父子进程`模式以及`Master-Worker`模式。

### v0.0.8 (2021-10-22)

1. 完善`平滑重启`的逻辑。
2. `Supervisor进程监控管理`模块新增接口，支持从`supervisor.ini`格式的配置文件中载入配置，从而启动配置的进程。

### v0.0.7

1. 完成`EventBus(事件总线)`功能的开发，详情请见 [eventBus(事件总线)](component/eventBus.md)
2. 修复`Supervisor进程监控管理`模块的一些bug。

### v0.0.6

1. 增加`Supervisor进程监控管理`模块，初步支持管理多个子进程，后续会支持通过`supervisor.conf`配置文件配置子进程启动。
2. 完善`EasyService`快速服务创建模块。
3. 完善`drpc message`模块的测试用例及文档。

### v0.0.5

1. 增加 `BeforeCloseEndpoint` 关闭Endpoint之前触发该事件.
2. 增加 `AfterCloseEndpointPlugin` 关闭Endpoint之后触发该事件.
3. 完善 `event`事件的测试用例。
4. 完善 `confg`配置的测试用例。