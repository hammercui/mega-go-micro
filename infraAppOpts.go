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

type AppOpts struct {
	IsConfWatchOn  bool
	IsBrokerOn     bool
	IsRedisOn      bool
	IsMongoOn      bool
	IsSqlOn        bool
	IsSkyWalkingOn bool
}

func InitAppWithOpts(opts *AppOpts) *InfraApp {
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
	app = InfraApp{
		Reg:      reg,
		Selector: sel,
	}

	//4
	if opts.IsSqlOn {
		app.ReadOnlyDB = mysql.NewMysqlReadOnly()
		app.ReadWriteDB = mysql.NewMysqlReadWrite()
	}

	//5 新建配置中心合并配置
	if opts.IsConfWatchOn {
		app.ConfWatch = watch.NewConfWatch()
	}

	//6 redis client
	if opts.IsRedisOn {
		app.RedisClient = infraRedis.NewRedisClient()
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
	if opts.IsConfWatchOn{
		regisConfWatch()
	}

	return &app
}
