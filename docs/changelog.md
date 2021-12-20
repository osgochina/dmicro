## 更新日志

### v0.2.0
1. 支持`ProtoBuf`协议.
2. 增加安全传输`SecureBodyPlugin`插件.
3. 完善文档，增加框架`logo`.

### v0.1.4
1. MacOS 支持
2. Windows 支持
3. `平滑重启` 逻辑支持，支持`父子进程`模式以及`Master-Worker`模式。

### v0.0.8

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