# mega-go-micro
>基于go-micro v2封装的golang微服务框架,使用gin作为http路由功能

## 基础组件服务
>提供mysql,redis,mongo.consul,kafka等服务的连接池
>提供工具函数

安装
```
go get github.com/hammercui/mega-go-micro
```
support flag 
* -configs `eg: -configs=/data/docker/micro/configs`
* -version 
* -env `eg:dev,coder,bea,prod`
* -nodeId `eg: 1`
* -ip: `eg: 182.168.1.10`

## How to use

### FAQ

1 etcd err:
go mod 使用gprc 1.26.0
```
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
```
## 更新日志
### v2.1.0路线图
* [ ] infraApp变更为facade模式
### v2.0.0路线图
* [x] 重构configs配置系统,修改为根据env读取配置文件,比如application.yml,application-dev.yml
* [x] mysql多数据库支持
* [x] redis多数据库支持
* [x] 日志支持配置存活时间,避免磁盘溢出


### v1.2.7
* 增加读写分离支持
* 修复redis无密码时连接错误