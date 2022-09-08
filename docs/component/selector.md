# 服务发现 (Selector)

## 概念

`RPC Server` 服务通过 `Registry` 功能注册后，它的调用者`RPC Client`需要能方便的发现该服务，并请求它。
这是最基本的功能，当该服务中的某个节点下线后，它的调用者要能够及时的剔除该节点。

针对`节点`的情况，使用 `Strategy` 策略引擎分配给节点不同的权重，
根据节点信息使用 `Filter` 过滤器剔除节点。

> ps: `Selector` 是提供给 `RPC Client`使用。需要配合 `Registry` 一起使用。


## 功能特性

* 方便的发现服务，选择可用节点。
* 使用 `Strategy` 策略引擎分配给节点不同的权重。
* 使用 `Filter` 过滤器剔除节点。
* 实时感知节点状态。

## 快速使用

### 使用默认的 `MDNS` Registry组件

关于`mdns`的介绍，可以查看[mdns](https://en.wikipedia.org/wiki/Multicast_DNS)

```go
cli := client.NewRpcClient(serviceName,
    client.OptSelector(
        selector.NewSelector(
            selector.OptRegistry(registry.DefaultRegistry),
        ),
    ),
)
```

`registry.DefaultRegistry` 默认的服务注册中心使用的是`mdns`,这样你就可以在本地局域网直接通过`rpc client`请求到该服务。

### 使用 `Etcd` Registry组件

```go
1. serviceName := "testregistry"
2. etcd.SetPrefix("/vprix/registry/dev/")
3. cli := client.NewRpcClient(serviceName,
4.     client.OptSelector(
5.         selector.NewSelector(
6.             selector.OptRegistry(
7.                 etcd.NewRegistry(
8.                     registry.OptAddrList("127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"),
9.                     registry.OptServiceName(serviceName),
10.                 ),
11.             ),
12.         ),
13.     ),
14. )
```

* 第2行，设置服务在`etcd`中的存储路径。
* 第8行，设置`etcd`集群的地址.


### 使用 `Memory` Registry组件

```go
svr := &registry.Service{
    Nodes: []*registry.Node{
        {
            Address: "127.0.0.1:9091",
        },
    },
}
cli := client.NewRpcClient("testregistry", client.OptCustomService(svr))
```

通过编码中显示的设置服务节点地址，请求服务。

