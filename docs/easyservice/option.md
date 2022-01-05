# 启动参数

```shell
USAGE
	./server [start|stop|quit] [default|custom] [OPTION]
OPTION
	-c,--config     指定要载入的配置文件，该参数与gf.gcfg.file参数二选一，建议使用该参数
	-d,--daemon     使用守护进程模式启动
	--env           环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product
	--debug         是否开启debug 默认debug=false
	--pid           设置pid文件的地址，默认是/tmp/[server].pid
	-h,--host       服务监听地址
	-p,--port       服务监听端口
	--network       监听的网络协议，支持tcp,tcp4,tcp6,默认tcp
	
EXAMPLES
	/path/to/server 
	/path/to/server start --env=dev --debug=true --pid=/tmp/server.pid
	/path/to/server start --host=127.0.0.1 --port=8808
	/path/to/server start --host=127.0.0.1 --port=8808 --gf.gcfg.file=config.product.toml
	/path/to/server start --host=127.0.0.1 --port=8808 -c=config.product.toml
	/path/to/server start user --host=127.0.0.1 --port=8808 
	/path/to/server start pay,user,admin  --host=127.0.0.1 --port=8808
	/path/to/server stop
	/path/to/server quit
	/path/to/server reload
	/path/to/server version
	/path/to/server help
```

### 启动命令

#### start

启动服务。

如果启动程序的时候，不传入参数，则默认执行`start`参数的逻辑。
第二个参数为要启动的服务名，可以传入多个，以逗号分割。也可以不传入，执行默认逻辑。
具体的服务启动对应关系，需要开发者自行掌握，通过`SandboxNames`方法获取传入的参数。

```shell
$ /path/to/server 
$ /path/to/server start 
$ /path/to/server start pay
$ /path/to/server start pay,user,admin
```
以上四种方式都是支持的。

#### stop

强制停止服务。

请注意，执行停止命令的时候，请保持`start`时候一样的参数。原因是会使用到对应的pid文件发送信号。
如：
* `/path/to/server ` 对应 `/path/to/server stop`
* `/path/to/server start ` 对应 `/path/to/server stop`
* `/path/to/server start pay` 对应 `/path/to/server stop pay`
* `/path/to/server start pay,user,admin` 对应 `/path/to/server stop pay,user,admin`

该命令不会等待服务处理完毕，进程收到信号后直接退出。

#### quit

优雅的退出服务。

请注意，执行退出命令的时候，请保持`start`时候一样的参数。原因是会使用到对应的pid文件发送信号。
如：
* `/path/to/server ` 对应 `/path/to/server quit`
* `/path/to/server start ` 对应 `/path/to/server quit`
* `/path/to/server start pay` 对应 `/path/to/server quit pay`
* `/path/to/server start pay,user,admin` 对应 `/path/to/server quit pay,user,admin`

该命令会等待服务处理完现有的工作，再退出，最多等待30秒。

#### reload

平滑重启服务。

请注意，执行重启命令的时候，请保持`start`时候一样的参数。原因是会使用到对应的pid文件发送信号。
如：
* `/path/to/server ` 对应 `/path/to/server reload`
* `/path/to/server start ` 对应 `/path/to/server reload`
* `/path/to/server start pay` 对应 `/path/to/server reload pay`
* `/path/to/server start pay,user,admin` 对应 `/path/to/server reload pay,user,admin`

该命令支持服务的平滑重启，具体的平滑重启使用逻辑，请参考[平滑重启](drpc/graceful.md).

#### version

显示当前程序的编译版本。

#### help

展示帮助信息。

### 各项参数说明

#### `-c,--config`

指定要载入的配置文件

启动服务的时候，针对配置文件，默认规则如下：
1. 有`-c,--config`参数，则优先使用该参数。
2. 未传入`-c,--config`参数，而是传入了`--gf.gcfg.file`,`--gf.gcfg.path`,则使用gf框架的配置获取流程，注意该参数可以使用环境变量传入`GF_GCFG_FILE`,`GF_GCFG_PATH`。
3. 如果以上两个参数都未传入，则为了屏蔽gf框架自带的报错，自动生成一段空的配置，使用`g.Cfg().Get()`等方法获取配置返回为空。

4. 如果参数传入的是文件名，则会去`./`,`./config`,`mainPkg/`,`mainPkg/config/`去查找文件。
```shell
	$ /path/to/server -c=config.toml
```
5. 如果参数传入的是文件路径，则会去该路径读取配置。如：

```shell
	$ /path/to/server -c=./config.toml
	$ /path/to/server -c=/path/to/config.toml
```

#### `--gf.gcfg.file`
设置配置文件名
获取配置相关的参数还有：
* gf.gcfg.path
* gf.gcfg.errorprint

参考gf的文档，https://goframe.org/pages/viewpage.action?pageId=1114668

#### `-d,--daemon`

使用守护进程模式启动

#### `--env`

环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product

#### `--debug`

是否开启debug 默认debug=false

#### `--pid`

设置pid文件的地址，默认是/tmp/[server].pid

#### `-h,--host`

服务监听地址,该参数适用于使用`BoxConf`获取配置的模式，如果有多个sandbox，可以通过配置文件获取。

#### `-p,--port`

服务监听端口,该参数适用于使用`BoxConf`获取配置的模式，如果有多个sandbox，可以通过配置文件获取。

#### `--network `

监听的网络协议，支持tcp,tcp4,tcp6,unixsocket,kcp,quic,默认tcp,该参数适用于使用`BoxConf`获取配置的模式，如果有多个sandbox，可以通过配置文件获取。




