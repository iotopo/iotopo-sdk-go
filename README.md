# iotopo-sdk-go

## 缓存
* 普通 redis 缓存
* 对象缓存，基于 redis 实现，支持本地二级缓存

## 分布式锁
基于 redis 实现

## RPC
基于 nrpc 实现, 依赖 nats 中间件：
https://github.com/nats-rpc/nrpc

nRPC needs Go 1.7 or higher. $GOPATH/bin needs to be in $PATH for the protoc invocation to work. To generate code, you need the protobuf compiler (which you can install from here) and the nRPC protoc plugin.

To install the nRPC protoc plugin:
```shell script
$ go get github.com/nats-rpc/nrpc/protoc-gen-nrpc
```

## 消息订阅发布
基于 nats 实现, 依赖 nats 中间件