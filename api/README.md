# API

API 基于 1uvu/serve 和 1uvu/fabric-sdk-client 实现了一个用于访问 Fabric 网络的 API Server, 可以使用简洁的 HTTP API 接口来执行**通用**的链码调用 (已完成), 交易查询等操作, 因此它可以适应各种 Fabric 网络, 只需满足按照规定的格式请求即可.

- config: 配置各个组织的身份信息, 用于生成 Client, 其中 Client 分为 app 和 admin 两种, 可按需要配置
- client: 基于 config 提供 Client 单例, 调用 GetApp 或 GetAdmin 可并发地获取 Client 单例
- server: 定义接口路由和请求处理逻辑, 返回接口请求结果, 执行 Run 来启动 API Server

## 调用格式

1. url

   | field | meaning |
   |  ----  | ----  |
   | :group | 使用何种 Client (app or admin)app/org1/channel1/chaincode/invoke |
   | :orgid | 组织 id |
   | :channelid | 通道 id |
   | :client | 使用哪种具体的 Client 实例 (包括 chaincode, ledger, msp, resource 四种, 详见代码注释) |
   | :opt | 每种 Client 实例支持的操作 |

    当前 client 只支持设为 chaincode, opt 只支持 invoke

2. post body

   | field | meaning |
   |  ----  | ----  |
   | chaincodeID | 链码 id (string) |
   | fcn | 合约名称 (string) |
   | args | 合约参数列表 (string 列表) |
   | needSubmit | 是否提交到账本, 即写入区块链 (true or false) |
   | endpoints | 调用节点的 Hosts 列表 (需要遵循背书策略, 建议设为空列表, 使用默认) |

3. reponse body

   | field | meaning |
   |  ----  | ----  |
   | payload | 返回调用结果的 base64 编码 |
   | transactionInfo | 交易信息 (只针对 admin Client), 包括交易 id, 交易发起组织, 以及交易时间|
   | chaincodeStatus | 链码调用响应状态码 (只针对 admin Client), 调用成功则为 200 |

## 调用示例

如在如下场景中:
> 组织 Org1 在 Channel1 中部署了名为 patient 的链码, 其中提供了 Register 和 Query 两个合约, 合约签名 (简化版本, 具体链码可查看 [channel1/patient](../chaincode/channel1/patient/patient.go))分别为:

- Register(patientID string, patient Patient)
- Query(patientID string) Patient

那么可以提供如下方式调用:

1. 使用 curl:

   ```shell

   # 注意链接的参数格式

   curl -X POST -d {\"chaincodeID\":\"patient\"\,\"fcn\":\"Register\"\,\"args\":[\"p1\"\,\"{\'name\':\'ZJH-1\'\,\'gender\':\'male\'\,\'birth\':\'2000-10-01\'\,\'identifyID\':\'xxxxxx-xxxx-20001001-xxxx-xxxx\'\,\'phoneNumber\':\'111-2200-0000\'\,\'address\':\'CQ\'\,\'nativePlace\':\'NG\'\,\'creditCard\':\'6217-0000-0000-0000\'\,\'healthcareID\':\'h1\'}\"]\,\"needSubmit\":true\,\"endpoints\":[]} http://127.0.0.1:9999/admin/org1/channel1/chaincode/invoke


   curl -X POST -d {\"chaincodeID\":\"patient\"\,\"fcn\":\"Query\"\,\"args\":[\"p1\"]\,\"needSubmit\":false\,\"endpoints\":[]} http://127.0.0.1:9999/app/org1/channel1/chaincode/invoke

   ```

2. 使用 REST Client

   | url | post body |
   |  ----  | ----  |
   | http://127.0.0.1:9999/admin/org1/channel1/chaincode/invoke | {"chaincodeID":"patient","fcn":"Register","args":["p1","{'name':'ZJH-1','gender':'male','birth':'2000-10-01','identifyID':'xxxxxx-xxxx-20001001-xxxx-xxxx','phoneNumber':'111-2200-0000','address':'CQ','nativePlace':'NG','creditCard':'6217-0000-0000-0000','healthcareID':'h1'}"],"needSubmit":true,"endpoints":[]} |
   | http://127.0.0.1:9999/app/org1/channel1/chaincode/invoke | {"chaincodeID":"patient","fcn":"Query","args":["p1"],"needSubmit":"false","endpoints":[]} |

