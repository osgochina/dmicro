# Ctx 上下文

在`drpc`服务中，每个请求的`request`都是单独开启一个`goroutine`进行处理，`Ctx`对象中包含了该次请求的所有信息。请求结束以后会被销毁。



## 如何使用

### `CallCtx`

## `Ctx`支持的的方法

<table>
<thead>
<td>方法名</td>
<td>功能</td>
<td>EarlyCtx</td>
<td>WriteCtx</td>
<td>ReadCtx</td>
<td>PushCtx</td>
<td>CallCtx</td>
<td>UnknownPushCtx</td>
<td>UnknownCallCtx</td>
</thead>
<tr>
  <td>Endpoint()</td>
  <td>获取当前endpoint对象</td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray" ></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td  style="background-color:gray"></td>
  <td  style="background-color:gray"></td>
  <td  style="background-color:gray"></td>
</tr>
<tr>
  <td>Session()</td>
  <td>获取当前会话对象</td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td  style="background-color:gray"></td>
</tr>
<tr>
  <td>IP()</td>
  <td>返回远端ip</td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
</tr>
<tr>
  <td>RealIP()</td>
  <td>返回远端真实ip</td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
</tr>
<tr>
  <td>Swap()</td>
  <td>返回自定义交换区数据</td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
</tr>
<tr>
  <td>Context()</td>
  <td>获取上下文</td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
  <td style="background-color:gray"></td>
</tr>
<tr>
  <td>Input()</td>
  <td>获取传入的消息</td>
  <td></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td></td>
  <td></td>
</tr>
<tr>
  <td>Output()</td>
  <td>将要发送的消息对象</td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td ></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td ></td>
  <td ></td>
</tr>
<tr>
  <td>StatusOK()</td>
  <td>状态是否ok</td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td style="background-color: gray"></td>
  <td ></td>
  <td ></td>
  <td ></td>
  <td ></td>
</tr>
<tr>
  <td>Status()</td>
  <td>当前消息状态</td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td style="background-color: gray"></td>
  <td ></td>
  <td ></td>
  <td ></td>
  <td ></td>
</tr>
<tr>
  <td>Seq()</td>
  <td>获取消息的序列号</td>
  <td ></td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>PeekMeta()</td>
  <td>窥视消息的元数据</td>
 <td ></td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>VisitMeta()</td>
  <td>浏览消息的元数据</td>
 <td ></td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>CopyMeta()</td>
  <td>获取消息的元数据副本</td>
 <td ></td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>ServiceMethod()</td>
  <td>该消息需要访问的服务名</td>
 <td ></td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>ResetServiceMethod()</td>
  <td>重置该消息将要访问的服务名</td>
 <td ></td>
  <td ></td>
  <td style="background-color:gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>GetBodyCodec()</td>
  <td>获取当前消息的编码格式</td>
  <td></td>
  <td></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>InputBodyBytes()</td>
  <td>传入的消息[]byte数组</td>
  <td></td>
  <td></td>
  <td></td>
  <td></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>BuildBody()</td>
  <td>构建消息内容</td>
  <td></td>
  <td></td>
  <td></td>
  <td></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>SetBodyCodec()</td>
  <td>设置回复消息的编码格式</td>
  <td></td>
  <td></td>
  <td></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td></td>
  <td  style="background-color: gray"></td>
</tr>
<tr>
  <td>SetMeta()</td>
  <td>设置指定key的值</td>
  <td></td>
  <td></td>
  <td></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td></td>
  <td style="background-color: gray"></td>
</tr>
<tr>
  <td>AddTFilterId()</td>
  <td>设置回复消息传输层的编码过滤方法id</td>
  <td></td>
  <td></td>
  <td></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td></td>
  <td style="background-color: gray"></td>
</tr>

<tr>
  <td>ReplyBodyCodec()</td>
  <td>获取响应消息的编码格式</td>
  <td></td>
  <td></td>
  <td></td>
  <td></td>
  <td style="background-color: gray"></td>
  <td></td>
  <td></td>
</tr>
</table>

## 不同Ctx的应用场景

`Ctx`根据服务的不同阶段，其内容和含义也不一样。

业务开发的时候能用到四种`Ctx`:

- `CallCtx` 注册Call方法时使用，提供`Call`方法被`endpoint`调用(Call方法有返回值).
- `PushCtx` 注册Push方法时使用，提供`Push`方法被`endpoint`调用(Push方法没有返回值)。
- `UnknownCallCtx` 如果请求Call方法时候为命中，则统一路由到某个方法，其参数为`UnknownCallCtx`。
- `UnknownPushCtx` 如果请求Push方法时候为命中，则统一路由到某个方法, 其参数为`UnknownPushCtx`。

还有几种`Ctx`，正常的业务开发不太用到，但是在进行插件开发的时候，根据请求的不同生命阶段，会用到它们。
当然，插件开发不仅限于以下三种，而是根据hook节点不一样可能会用到所有`Ctx`.

- `EarlyCtx` 

   用在`BeforeReadHeader`hook阶段，执行读取Header之前触发该事件，作为参数传入。

- `WriteCtx` 

    所有消息的`写`相关hook阶段，作为参数传入，它们分别是：
  - BeforeWriteCall        写入 CALL 消息之前执行该事件
  - AfterWriteCall         写入CALL消息之后执行该事件
  - BeforeWriteReply       写入REPLY消息之前执行该事件
  - AfterWriteReply        写入Reply消息之后执行该事件
  - BeforeWritePush        写入PUSH消息之前执行该事件
  - AfterWritePush         写入PUSH消息之后执行该事件
  
- `ReadCtx` 

    所有消息的`读`相关hook阶段，作为参数传入，它们分别是：
  - AfterReadCallHeader    读取CALL消息的Header之后触发该事件
  - BeforeReadCallBody     读取CALL消息的body之前触发该事件
  - AfterReadCallBody      读取CALL体之后执行该事件
  - AfterReadPushHeader    读取PUSH消息头之后执行该事件
  - BeforeReadPushBody     读取PUSH消息体之前执行该事件
  - AfterReadPushBody      读取PUSH消息体之后执行该事件
  - AfterReadReplyHeader   读取Reply消息头之后执行该事件
  - BeforeReadReplyBody    读取Reply消息体之前执行该事件
  - AfterReadReplyBody     读取Reply消息体之后执行该事件



## 实现原理

