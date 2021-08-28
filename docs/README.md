## dmicro简介

> dmicro是一个高效、可扩展且简单易用的微服务框架。包含drpc,easyserver等组件。

该项目的诞生离不开`erpc`和`GoFrame`两个优秀的项目。

其中`drpc`组件参考`erpc`项目的架构思想，依赖的基础库是`GoFrame`。

* [erpc](https://gitee.com/henrylee/erpc)
* [GoFrame](https://gitee.com/johng/gf)


## 安装

```html
go get -u -v github.com/osgochina/dmicro
```
推荐使用 `go.mod`:
```
require github.com/osgochina/dmicro latest
```

* import

```go
import "github.com/henrylee2cn/erpc/v6"
```

## 限制
```shell
golang版本 >= 1.13
```
