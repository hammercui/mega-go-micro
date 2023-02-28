package base

import (
	infraBroker "github.com/hammercui/mega-go-micro/v2/broker"
	"github.com/hammercui/mega-go-micro/v2/conf"
	"github.com/hammercui/mega-go-micro/v2/log"
	"github.com/hammercui/mega-go-micro/v2/mongo"
	"github.com/hammercui/mega-go-micro/v2/mysql"
	infraRedis "github.com/hammercui/mega-go-micro/v2/redis"
	"github.com/hammercui/mega-go-micro/v2/registry/consul"
	"github.com/hammercui/mega-go-micro/v2/tracer"
	"github.com/hammercui/mega-go-micro/v2/watch"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
)

// 初始化app
func InitApp() *InfraApp {
	//1 配置初始化
	conf.InitConfig()
	_conf := conf.GetConf()
	//2 日志初始化
	log.InitLog()
	log.Logger().Infof("-------log init over-------")
	//3 consul注册
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = _conf.Consul.Addrs
	})
	sel := selector.NewSelector(selector.Registry(reg))

	//4 TODO 配置中心合并配置，使用外部注入的方式
	confWatch := watch.InitConfWatch()
	//5 init redis
	redisMap := infraRedis.InitRedis()
	log.Logger().Infof("-------redis init over-------")
	//6 init broker
	brokerIns := infraBroker.InitKafkaBroker()
	log.Logger().Info("-------kafka init over-------")
	//7 init trace
	tracerIns := tracer.InitTracer()
	log.Logger().Infof("-------tracer init over-------")
	//8 init mysql
	mysqlMap := mysql.InitMysql()
	log.Logger().Infof("-------mysql init over-------")
	//9 init mongo
	mongoMap := mongo.InitMongo()
	log.Logger().Infof("-------mongo init over-------")

	// 初始化
	app = &InfraApp{
		Reg:       reg,
		Selector:  sel,
		RedisMap:  redisMap,
		Broker:    brokerIns,
		ConfWatch: confWatch,
		Tracer:    tracerIns,
		MySqlMap:  mysqlMap,
		MongoMap:  mongoMap,
	}

	//9 监听配置
	//regisConfWatch()
	return app
}

// 卸载app
func UnitApp() {
	//clear mongo client
}
