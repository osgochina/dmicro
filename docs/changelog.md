## 更新日志

### v0.0.6

1. 增加`Supervisor进程监控管理`模块，初步支持管理多个子进程，后续会支持通过`supervisor.conf`配置文件配置子进程启动。
2. 完善`EasyService`快速服务创建模块。
3. 完善`drpc message`模块的测试用例及文档。

### v0.0.5

1. 增加 `BeforeCloseEndpoint` 关闭Endpoint之前触发该事件.
2. 增加 `AfterCloseEndpointPlugin` 关闭Endpoint之后触发该事件.
3. 完善 `event`事件的测试用例。
4. 完善 `confg`配置的测试用例。