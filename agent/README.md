# agent(网关）

[![Build Status](https://travis-ci.org/gonet2/agent.svg?branch=master)](https://travis-ci.org/gonet2/agent)

## 特性

1. 处理各种协议的接入，同时支持TCP和UDP(KCP协议)，进行双栈通信。
1. 连接管理，会话建立，数据包加解密(DH+RC4)。
1. **透传**解密后的原始数据流到后端（通过gRPC streaming)。
1. **复用**多路用户连接，到一条通往game的物理连接。
1. 不断开连接切换后端业务。
1. 唯一入口，安全隔离核心服务。

## 协议号划分

数据包会根据协议编号（0-65535) **透传** 到对应的服务， 例如(示范）:      

      1-1000: 登陆相关协议，网关协同auth服务处理。
      2001-3000: 游戏逻辑段
      ....
      
[查看消息协议号定义](http://gitlab.airdroid.com/chess/chess/blob/dev/grpc/agent.proto)

## 消息封包格式
 
        +----------------------------------------------------------------+     
        | SIZE(2) | TIMESTAMP(4) | PROTO(2) | PAYLOAD(SIZE-6)            |     
        +----------------------------------------------------------------+     
        
> SIZE: 后续数据包总长度         
> TIMESTAMP: 数据包序号           
> PROTO: 协议号           
> PAYLOAD: 负载           
