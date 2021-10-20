package infra

import (
	"fmt"
	infraBroker "github.com/hammercui/mega-go-micro/broker"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"github.com/hammercui/mega-go-micro/mysql"
	infraRedis "github.com/hammercui/mega-go-micro/redis"
)

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
		app.ReadWriteDB = mysql.DefaultMysqlReadWrite()
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
		app.ReadOnlyDB = mysql.DefaultMysqlReadOnly()
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
		app.RedisClient = infraRedis.DefaultRedisClient()
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