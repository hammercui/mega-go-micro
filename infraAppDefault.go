package infra

import (
	infraBroker "github.com/hammercui/mega-go-micro/broker"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"github.com/hammercui/mega-go-micro/mysql"
	infraRedis "github.com/hammercui/mega-go-micro/redis"
	"github.com/hammercui/mega-go-micro/registry/consul"
	skyWalking2 "github.com/hammercui/mega-go-micro/tracer/skyWalking"
	"github.com/hammercui/mega-go-micro/watch"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
)

/**
infraApp 默认方式生成
*/

//初始化app
func InitApp() *InfraApp {
	//1 配置初始化
	conf.InitConfig()
	//2 日志初始化
	log.InitLog(nil)
	//3 自定义consul注册
	consulConf := conf.GetConf().ConsulConf
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = consulConf.Addrs
	})
	sel := selector.NewSelector(selector.Registry(reg))

	//4 新建配置中心合并配置
	confWatch := watch.NewConfWatch()
	//5 redis client
	redisClient := infraRedis.DefaultRedisClient()
	//6 init broker
	brokerIns := infraBroker.NewKafkaBroker()

	skyWalking := skyWalking2.NewSkyTracer()

	//7 初始化
	app = &InfraApp{
		ReadOnlyDB:         mysql.DefaultMysqlReadOnly(),
		ReadWriteDB:        mysql.DefaultMysqlReadWrite(),
		Reg:                reg,
		Selector:           sel,
		RedisClient:        redisClient,
		Broker:             brokerIns,
		ConfWatch:          confWatch,
		SkyWalking:         skyWalking,
		readOnlyDBPoolMap:  make(map[string]*mysql.DBPool),
		readWriteDBPoolMap: make(map[string]*mysql.DBPool),
		redisPoolMap:       make(map[string]*infraRedis.RedisPool),
	}
	//8 池化
	appConf := conf.GetConf()
	app.SetReadOnlyDBPool(appConf.AppConf.Name, mysql.DefaultMysqlDsn(), app.ReadOnlyDB)
	app.SetReadWriteDBPool(appConf.AppConf.Name, mysql.DefaultMysqlDsn(), app.ReadWriteDB)
	app.SetRedisPool(appConf.AppConf.Name,appConf.RedisConf.Addr,appConf.RedisConf.DbIndex,app.RedisClient)
	//9 监听配置
	regisConfWatch()

	return app
}
