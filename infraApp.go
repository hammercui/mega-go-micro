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
	infraBroker "github.com/hammercui/mega-go-micro/broker"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"github.com/hammercui/mega-go-micro/mysql"
	infraRedis "github.com/hammercui/mega-go-micro/redis"
	"github.com/hammercui/mega-go-micro/registry/consul"
	"github.com/hammercui/mega-go-micro/watch"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"gorm.io/gorm"
)

type InfraApp struct {
	HttpRunning bool
	RpcRunning  bool
	RedisClient *redis.Client
	ReadOnlyDB  *gorm.DB
	ReadWriteDB *gorm.DB
	Reg         registry.Registry //服务注册
	Selector    selector.Selector //服务发现
	Broker      broker.Broker     //消息订阅与发布
	//配置中心
	ConfWatch *watch.ConfWatch
	//todo mongo连接
}

var app InfraApp

//初始化app
func InitApp() *InfraApp {
	//1 配置初始化
	conf.InitConfig()
	//2 日志初始化
	log.InitLog()
	//3 自定义consul注册
	consulConf := conf.GetConf().ConsulConf
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = consulConf.Addrs
	})
	sel := selector.NewSelector(selector.Registry(reg))

	//4 从配置中心合并配置并监听配置
	confWatch := watch.InitConfWatch()

	//5 redis client
	redisClient := infraRedis.InitRedis()
	//6 init broker
	brokerIns := infraBroker.NewKafkaBroker()

	//初始化
	app = InfraApp{
		ReadOnlyDB:  mysql.InitMysqlReadOnly(),
		ReadWriteDB: mysql.InitMysqlReadWrite(),
		Reg:         reg,
		Selector:    sel,
		RedisClient: redisClient,
		Broker:      brokerIns,
		ConfWatch:   confWatch,
	}

	regisConfWatch()

	return &app
}

func regisConfWatch() {
	//mysql
	app.ConfWatch.Watch("mysql", &map[string]string{}, func(outConf interface{}, err error) {
		log.Logger().Debugf("监听到配置变更", outConf)
		//todo mysql重连
	})
}

//卸载app
func UnitApp() {

}
