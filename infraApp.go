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
	"fmt"
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

	//4 新建配置中心合并配置
	confWatch := watch.NewConfWatch()
	//5 redis client
	redisClient := infraRedis.NewRedisClient()
	//6 init broker
	brokerIns := infraBroker.NewKafkaBroker()

	//7 初始化
	app = InfraApp{
		ReadOnlyDB:  mysql.NewMysqlReadOnly(),
		ReadWriteDB: mysql.NewMysqlReadWrite(),
		Reg:         reg,
		Selector:    sel,
		RedisClient: redisClient,
		Broker:      brokerIns,
		ConfWatch:   confWatch,
	}
	//8 监听配置
	regisConfWatch()

	return &app
}

func regisConfWatch() {
	//mysql
	app.ConfWatch.Watch("mysql", &map[string]string{}, func(outConf interface{}, err error) {
		if err != nil {
			return
		}
		log.Logger().Info("trigger mysql config change: ", outConf)
		//mysql重连
		mysqlMap := outConf.(*map[string]string)
		conf.GetConf().MysqlConf.Addr = fmt.Sprintf("%s:%s", (*mysqlMap)["host"], (*mysqlMap)["port"])
		app.ReadWriteDB = mysql.NewMysqlReadWrite()
		log.Logger().Info("trigger mysql reconnect success!")
	})

	//readMysql
	app.ConfWatch.Watch("readMysql", &map[string]string{}, func(outConf interface{}, err error) {
		if err != nil {
			return
		}
		log.Logger().Info("trigger readMysql config change: ", outConf)
		//mysql重连
		readMysqlMap := outConf.(*map[string]string)
		conf.GetConf().MysqlConf.ReadAddr = fmt.Sprintf("%s:%s", (*readMysqlMap)["host"], (*readMysqlMap)["port"])
		app.ReadOnlyDB = mysql.NewMysqlReadOnly()
		log.Logger().Info("trigger readMysql reconnect success!")
	})

	//redis
	app.ConfWatch.Watch("redis", &[]map[string]interface{}{}, func(outConf interface{}, err error) {
		if err != nil {
			return
		}
		log.Logger().Info("trigger redis config change: ", outConf)
		var redisMap = outConf.(*[]map[string]interface{})
		var redisAddrs []string
		for _, item := range *redisMap {
			redisAdds := fmt.Sprintf("%s:%v", item["host"], item["port"])
			redisAddrs = append(redisAddrs, redisAdds)
		}
		conf.GetConf().RedisConf.Sentinels = redisAddrs
		app.RedisClient.Close()
		app.RedisClient = infraRedis.NewRedisClient()
		log.Logger().Info("trigger redis reconnect success!")
	})

	// kafka
	app.ConfWatch.Watch("kafka", &[]map[string]interface{}{}, func(outConf interface{}, err error) {
		if err != nil {
			return
		}
		log.Logger().Info("trigger kafka config change: ", outConf)
		var kafkaMap = outConf.(*[]map[string]interface{})
		var kafkaAddrs []string
		for _, item := range *kafkaMap {
			redisAdds := fmt.Sprintf("%s:%v", item["host"], item["port"])
			kafkaAddrs = append(kafkaAddrs, redisAdds)
		}
		conf.GetConf().KafkaConf.Addrs = kafkaAddrs
		app.Broker.Disconnect()
		app.Broker = infraBroker.NewKafkaBroker()
		log.Logger().Info("trigger kafka reconnect success!")
	})

	//todo mongo
	//app.ConfWatch.Watch("redis", &[]map[string]interface{}{}, func(outConf interface{}, err error) {
	//	if err != nil {
	//		return
	//	}
	//})
}

//卸载app
func UnitApp() {

}
