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
)

func init() {}

func InitRedis() map[string]*redis.Client {
	log.Logger().Infof("-------redis init console-------")
	_map := make(map[string]*redis.Client)
	for k, v := range conf.GetConf().RedisMap {
		if v.Enable {
			_map[k] = getClient(v)
		}
	}
	return _map
}

func getClient(redisConf *conf.RedisConf) *redis.Client {
	if redisConf.Sentinel == nil {
		return getClientByDirect(redisConf)
	} else {
		return getClientBySentinel(redisConf)
	}
}

//获得直连redis
func getClientByDirect(redisConf *conf.RedisConf) *redis.Client {
	opts := &redis.Options{
		Addr: redisConf.Addr,
		DB:   redisConf.DbIndex, // use default DB
	}
	if redisConf.Password != "" && len(redisConf.Password) > 1 {
		opts.Password = redisConf.Password
	}
	//connect redis
	redisClient := redis.NewClient(opts)
	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Logger().Errorf("redis direct connect fail! %s,err:%v", redisConf.Addr, err)
		panic(err)
	}
	log.Logger().Infof("redis direct connect success! %s %s !", redisConf.Addr, pong)
	return redisClient
}

//获得sentinel redis
func getClientBySentinel(redisConf *conf.RedisConf) *redis.Client {
	//connect redis
	flOpts := &redis.FailoverOptions{
		MasterName:    redisConf.Sentinel.Master,
		SentinelAddrs: redisConf.Sentinel.Nodes,
		DB:            redisConf.DbIndex,
	}
	if redisConf.Password != "" && len(redisConf.Password) > 1 {
		flOpts.Password = redisConf.Password
	}
	redisClient := redis.NewFailoverClient(flOpts)
	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Logger().Errorf("redis sentinel connect fail! %s,err:%v", redisConf.Sentinel.Nodes, err)
		panic(err)
	}
	log.Logger().Infof("redis sentinel connect success! %s %s ！", redisConf.Sentinel.Nodes, pong)
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
