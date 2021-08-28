# dmicro [![GitHub release](https://img.shields.io/github/release/osogchina/dmicro.svg?style=flat-square)](https://github.com/osgochina/dmicro/releases) [![report card](https://goreportcard.com/badge/github.com/osgochina/dmicro?style=flat-square)](http://goreportcard.com/report/osgochina/dmicro) [![github issues](https://img.shields.io/github/issues/osgochina/dmicro.svg?style=flat-square)](https://github.com/osgochina/dmicro/issues?q=is%3Aopen+is%3Aissue) [![github closed issues](https://img.shields.io/github/issues-closed-raw/osgochina/dmicro.svg?style=flat-square)](https://github.com/osgochina/dmicro/issues?q=is%3Aissue+is%3Aclosed) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/osgochina/dmicro) [![view examples](https://img.shields.io/badge/learn%20by-examples-00BCD4.svg?style=flat-square)](https://github.com/osgochina/dmicro/tree/main/.examples)

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

## 限制
```shell
golang版本 >= 1.13
```



