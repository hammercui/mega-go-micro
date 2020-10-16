/*
@Desc : 基础服务app对象
@Version : 1.0.0
@Time : 2020/8/22 16:02
@Author : hammercui
@File : infraApp
@Company: Sdbean
*/
package infra

import (
	"github.com/go-redis/redis"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/registry"
	"gorm.io/gorm"
	"wfServerMicro/infra/conf"
	"wfServerMicro/infra/log"
	"wfServerMicro/infra/mysql"
	infraRedis "wfServerMicro/infra/redis"
	infraBroker "wfServerMicro/infra/broker"
	"wfServerMicro/infra/registry/consul"
)

type InfraApp struct {
	HttpRunning bool
	RpcRunning  bool
	RedisClient *redis.Client
	ReadOnlyDB  *gorm.DB
	ReadWriteDB *gorm.DB
	Reg registry.Registry //服务注册与发现
	Broker broker.Broker
}

var app InfraApp

//初始化app
func InitApp() *InfraApp {
	//配置初始化
	conf.InitConfig()
	//日志初始化
	log.InitLog()
	//自定义consul注册
	consulConf := conf.GetConf().ConsulConf
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = consulConf.Addrs
	})
	//redis client
	redisClient := infraRedis.InitRedis()
	//init broker
	brokerIns := infraBroker.NewKafkaBroker()
	//初始化
	app = InfraApp{
		//RedisClient: infraRedis.InitRedis(),
		ReadOnlyDB: mysql.InitMysqlReadOnly(),
		ReadWriteDB:mysql.InitMysqlReadWrite(),
		Reg: reg,
		RedisClient: redisClient,
		Broker: brokerIns,
	}

	return &app
}

//卸载app
func UnitApp()  {
	
}
