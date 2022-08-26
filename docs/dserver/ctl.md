# 应用管理工具CTL

日常工作中，我们开发好应用发布运行后，只能通过日志查看其运行状态，或者通过外部命令，查看应用的cpu，内存占用情况，无法对其内部运行情况有更多的了解。

针对这种情况，`DServer`内置了ctl工具，让开发者很方便的就能够对开发的应用进行管理。

## 启动ctl

ctl工具默认开启，你也可以通过`dserver.CloseCtl()`方法关闭它。

1. 编译好应用后，启动该应用.
```shell
$ ./server
```
2. 打开第二个控制台窗口。
```shell
./server ctl

  ____    ____                                      
 |  _ \  / ___|    ___   _ __  __   __   ___   _ __ 
 | | | | \___ \   / _ \ | '__| \ \ / /  / _ \ | '__|
 | |_| |  ___) | |  __/ | |     \ V /  |  __/ | |   
 |____/  |____/   \___| |_|      \_/    \___| |_|  
Version:         No Version Info
Go Version:      No Version Info
DMicro Version:  v1.0.0
GF Version:      v1.16.9
Git Commit:      No Commit Info
Build Time:      No Time Info
Authors:         No Authors Info
Install Path:    /home/lzm/GolandProjects/dmicro/examples/dserver/server
DMicro »  
```
你进入到了一个全新的shell界面，可以在该界面管理你刚刚启动的服务。

## ctl管理

支持的命令

```shell
clear             clear the screen
debug             debug开关
exit              exit the shell
help              use 'help [command]' for command help
info, status, ps  查看当前服务状态
log               打印出服务的运行日志
reload            平滑重启服务
start             启动服务
stop              停止服务
version, v        打印当前程序的版本信息
```

1. 查看当前运行的服务 `info`,`status`,`ps`

    ```shell
    DMicro » ps
    ┌─────────────────┬─────────────┬─────────┬───────────────────────────┐
    │ SandBoxName     │ ServiceName │ Status  │ Description               │
    ├─────────────────┼─────────────┼─────────┼───────────────────────────┤
    │ DefaultSandBox  │ admin       │ Running │ pid 20911, uptime 0:03:14 │
    │ DefaultSandBox1 │ user        │ Running │ pid 20922, uptime 0:03:13 │
    └─────────────────┴─────────────┴─────────┴───────────────────────────┘
    DMicro »  
    ```

2. 开启关闭debug模式
    ```shell
    DMicro » debug open
    DMicro » debug close
    ```
   
3. 关闭服务`stop`
    ```shell
    DMicro » stop DefaultSandBox
    ``` 
4. 重启服务`reload`
    ```shell
    DMicro » reload DefaultSandBox
    ``` 
5. 启动服务`start`
    ```shell
    DMicro » start DefaultSandBox
    ``` 
6. 查看日志`log`
    ```shell
    DMicro » log --level info
    ``` 
   
> 注意，启动，关闭，重启命令的参数都是sandbox name，但是在多进程的模式下，操作的粒度是service。
也就是说重启某个sandbox，或重启该sandbox所在进程，其它在该进程的sandbox也会被重启。