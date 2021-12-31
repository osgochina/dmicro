## 编译参数的设置

`easyservice`支持特殊的编译变量。

* BuildVersion 指定的版本
* BuildGoVersion  编译时使用的go版本
* BuildGitCommitId 编译时git的commit id
* BuildTime        编译时间

如果在在编译的时候设置该变量，需要在编译时候加上参数：

```shell
go build -X "github.com/osgochina/dmicro/easyservice.BuildVersion=v0.1.1"
go build -X "github.com/osgochina/dmicro/easyservice.BuildGoVersion=go1.6.5"
go build -X "github.com/osgochina/dmicro/easyservice.BuildGitCommitId=fafsdffa"
go build -X "github.com/osgochina/dmicro/easyservice.BuildTime=2021-12-30"
```

## 使用编译脚本

编译脚本在`build/build.sh`。
如果大家需要使用，请把它复制到项目的根目录。

* 使用方式:
```build.sh [-s build_file] [-o output_dir] [-v version] [-g go_bin]```
* 参数详解:

    * build_file 需要从那个目录编译项目，默认是当前目录的main.go,你也可以指定指定的目录文件
    * output_dir 编译后的产物存在到那个目录.默认是存在在当前目录
    * version 编译后的文件版本号
    * go_bin 使用的golang程序

