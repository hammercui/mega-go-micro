/**
infraApp配置化生成
*/
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

func InitAppWithOpts(opts *conf.AppOpts) *InfraApp {
	//1 配置初始化
	conf.InitConfig()
	//2 日志初始化
	log.InitLog(opts)
	//3 自定义consul注册
	consulConf := conf.GetConf().ConsulConf
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = consulConf.Addrs
	})
	sel := selector.NewSelector(selector.Registry(reg))
	app = &InfraApp{
		Reg:                reg,
		Selector:           sel,
		readOnlyDBPoolMap:  make(map[string]*mysql.DBPool),
		readWriteDBPoolMap: make(map[string]*mysql.DBPool),
		redisPoolMap:       make(map[string]*infraRedis.RedisPool),
	}
	appConf := conf.GetConf()
	//4
	if opts.IsSqlOn {
		app.ReadOnlyDB = mysql.DefaultMysqlReadOnly()
		app.ReadWriteDB = mysql.DefaultMysqlReadWrite()
		//4.1 池化
		app.SetReadOnlyDBPool(appConf.AppConf.Name, mysql.DefaultMysqlReadOnlyDsn(), app.ReadOnlyDB)
		app.SetReadWriteDBPool(appConf.AppConf.Name, mysql.DefaultMysqlReadWriteDsn(), app.ReadWriteDB)
	}

	//5 新建配置中心合并配置
	if opts.IsConfWatchOn {
		app.ConfWatch = watch.NewConfWatch()
	}

	//6 redis client
	if opts.IsRedisOn {
		app.RedisClient = infraRedis.DefaultRedisClient()
		app.SetRedisPool(appConf.AppConf.Name, appConf.RedisConf.Addr, appConf.RedisConf.DbIndex, app.RedisClient)
	}

	//7 init broker
	if opts.IsBrokerOn {
		app.Broker = infraBroker.NewKafkaBroker()
	}
	//8 sky walking
	if opts.IsSkyWalkingOn {
		app.SkyWalking = skyWalking2.NewSkyTracer()
	}

	//8 监听配置
	if opts.IsConfWatchOn {
		regisConfWatch()
	}

	return app
}
