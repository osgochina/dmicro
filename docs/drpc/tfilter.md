# 传输管道过滤器(tfilter)
`tfilter`全称`transfer filter`,作用是对要进行网络传输的消息内容`pack`,对读取的消息内容执行`unpack`.

不同于`codec`,`tfilter`可以根据业务需求，按注册顺序多次执行。

`tfilter`位于`proto`与`codec`之间。

```mermaid
flowchart LR
     Message --> Proto --> TFilter --> Codec --> Socket
```

默认支持的传输管道过滤器(tfilter)

| id  | name     | 介绍                     |
|-----|----------|------------------------|
| m   | md5      | 对消息内容进行md5运算，保证消息的完整性  |
| g   | gzip     | 对消息内容进行zip压缩，提升消息的传输性能 |


## 如何使用传输管道过滤器



## 实现自己的`传输管道过滤器(tfilter)`