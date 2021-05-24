Fabric SDK 主要存在着五种客户端（客户端都是针对于当前组织来说的），其中常用的包括如下两种：
- Resource Manage Client (RC)
- Channel Client (CC)


其中 RC 用来管理 Fabric 网络中的各种资源，可以对 Fabric 网络进行各种操作，包括（源码角度）：
- pkg/fabsdk：主 package，主要用来生成 fabsdk 以及 fabric go sdk 中其他 pkg 使用的 option context。
- pkg/client/channel (CC)：主要用来调用、查询Fabric链码，或者注册链码事件。
- pkg/client/resmgmt (RC)：主要用来 Hyperledger fabric 网络的管理，比如创建通道、加入通道，安装、实例化和升级链码。
- pkg/client/event (EC)：配合 channel 模块来进行 Fabric 链码事件的注册和过滤。
- pkg/client/ledger (LC)：主要用来实现 Fabric 账本的查询，查询区块、交易、配置等。
- pkg/client/msp (MC)：主要用来管理 fabric 网络中的成员关系。

CC 可由 RC 获得，但是这对于分布在各地的应用来说，并不安全，因此 CC 还存在着另外一种方法来创建，对象类型虽然不同，但是其同样实现了 CC 和 EC 的全部功能，即：
- pkg/gateway

```golang
type Network struct {
	name    string
	gateway *Gateway
	client  *channel.Client
	event   *event.Client
}
```

使用 ccp config 生成 wallet，实现 cc 的功能（推荐）。

因此，这里推荐：
- 对于管理员客户端使用第一种方式来获取 cc。
- 对于普通用户应用，选用第二种方式来获取 cc。

所以，这里全部提供了

为了使用方便，这里将 cc 和 ec 封装为 AppClient，而另外 rc 和 mc 封装为 AdminClient。

如果想要使用 cc，ec 和 lc，则可以通过 admin client 的 sdk 成员十分简便地调用 admin.InitAppClient 实例化得到。

而且，并没有进行过多的封装，只是简化了实例化各种客户端的操作，扩展性更强。