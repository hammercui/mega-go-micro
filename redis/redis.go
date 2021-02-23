/*
@Desc : 微服务框架redis连接池
@Version : 1.0.0
@Time : 2020/8/25 14:37
@Author : hammercui
@File : redis
@Company: Sdbean
*/
package redis

import (
	//"fmt"
	"github.com/go-redis/redis"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"os"
)

//初始化redis
func initDirect() *redis.Client {
	redisConf := conf.GetConf().RedisConf
	//connect redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: "",                // no password set
		DB:       redisConf.DbIndex, // use default DB
	})
	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Logger().Infof("redis direct connect:%s fail,err:%v", redisConf.Addr, err)
		log.Logger().Error(err)
		os.Exit(0)
	}
	log.Logger().Infof("redis direct connect:%s %s !", redisConf.Addr, pong)
	return redisClient
}

func NewRedisClient() *redis.Client {
	appConf := conf.GetConf().AppConf
	redisConf := conf.GetConf().RedisConf
	if appConf.Env == conf.AppEnv_local {
		return initDirect()
	}

	//connect redis
	redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: redisConf.Sentinels,
		DB:            redisConf.DbIndex,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Logger().Error("redis sentinel connect fail!err:%v", err)
		os.Exit(0)
	}
	log.Logger().Info("redis sentinel connect success!%s ！", pong)
	return redisClient
}

//卸载redis
//func UnitRedis() {
//	if redisClient != nil {
//		//关闭redis连接
//		if err := redisClient.Close(); err != nil {
//			log.Logger().Error("redis close error:%v", err)
//		} else {
//			log.Logger().Info("redis close success!")
//		}
//	} else {
//		log.Logger().Info("redis is nil,no need close !")
//	}
//}
