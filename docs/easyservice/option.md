# 启动参数

```shell
USAGE
	./server [start|stop|quit] [default|custom] [OPTION]
OPTION
	-h,--host       服务监听地址，默认监听的地址为127.0.0.1
	-p,--port       服务监听端口，默认监听端口为0，表示随机监听
	--network       监听的网络协议，支持tcp,tcp4,tcp6,默认tcp
	-d,--daemon     debug模式开关，默认关闭debug=false
	--gf.gcfg.file  需要加载的配置文件名 如 config.dev.toml
	--env           环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product
	--debug         是否开启debug 默认debug=false
	--pid           设置pid文件的地址，默认是/tmp/[server].pid
EXAMPLES
	/path/to/server 
	/path/to/server start --env=dev --debug=true --pid=/tmp/server.pid
	/path/to/server start --host=127.0.0.1 --port=8808
	/path/to/server start --host=127.0.0.1 --port=8808 --gf.gcfg.file=config.product.toml
	/path/to/server start user --host=127.0.0.1 --port=8808 
	/path/to/server start pay  --host=127.0.0.1 --port=8808
	/path/to/server stop
	/path/to/server quit
	/path/to/server reload
	/path/to/server version
	/path/to/server help
```