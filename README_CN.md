# 阿里云云助手自建服务器

## 依赖
MySQL 用于持久化存储命令、执行结果、实例信息
Redis 频道订阅用于下发命令、临时存储执行结果


## 组件说明

* apiserver 用于与aliyun-assist-client进行交互
* manageserver 用于管理及下发命令


## 已实现功能

* RunCommand
* DesicribeInvocationResults
* 创建激活码
* 注册实例


## TODO

* [ ] manageserver 鉴权
* [ ] sendfile
* [ ] create command
* [ ] modify command
* [ ] delete command
* [ ] invoke command
* [ ] 实例在线状态维护