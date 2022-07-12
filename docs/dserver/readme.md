## DServer 服务管理功能

`dserver`服务管理功能能够让你专注与编写业务代码，编译部署后的运行时管理就交给它吧。
它支持单进程/多进程两种模式。

### 多进程模式

多进程模式分为`Master管理进程`,`业务处理进程`，该模式能够隔离业务，方便再一个仓库中开发多个业务，编译好以后能够通过一个服务，启动不同的业务进程。

`Master`进程管理者业务进程的整个声明周期，可以结束，拉起进程。也可以查看进程运行状态。

名词解释：

* DServer 整个项目中唯一，管理启动`master进程`
* Service 项目中可以有多个，每个Service是单独的进程。多个service之间不能通讯
* Sandbox 项目中可以有多个，并且每个service中可以有多个sandbox，同一个service中的sandbox可以相互通讯
* Graceful 平滑重启功能组件，运行于`DServer`的master进程。
* DServerCtl 服务管理客户端，支持链接到`master进程`，启动关闭服务，查看服务运行状态。

功能说明：

* 启动参数支持，支持按参数设置启动模式。
* `service进程`与`master进程`之间建立`unix socket`通讯。
* `master进程`暴露管理接口，支持通过`http`或者`rpc`方式调用。
* `DServerCtl`客户端支持应答式命令行。
* 各个进程的日志，支持独立写入，也支持转发到master进程统一处理。

待完成功能：

- [x] 通过反射注入sandbox需要的类型
- [ ] 增加统一的协程池，每个sandbox中的协程能统一管理
- [ ] 每个sandbox中有统一的Context上下文。
