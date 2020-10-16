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
	"os"
	"wfServerMicro/infra/conf"
	"wfServerMicro/infra/log"
)

var redisClient *redis.Client

//初始化redis
func initDirect() *redis.Client {
	redisConf := conf.GetConf().RedisConf
	//connect redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: "",                // no password set
		DB:       redisConf.DbIndex, // use default DB
	})
	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Logger().Infof("redis conn:%s fail,err:%v", redisConf.Addr, err)
		log.Logger().Error(err)
		os.Exit(0)
	}
	log.Logger().Infof("redis conn:%s %s !", redisConf.Addr, pong)
	return redisClient
}

func InitRedis() *redis.Client {
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
		log.Logger().Error("redis sentinel conn fail!err:%v", err)
		os.Exit(0)
	}
	log.Logger().Info("redis sentinel conn success!%s ！", pong)
	return redisClient

	////connect mongo
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//p.mongoClient,err = mongo.Connect(ctx,options.Client().ApplyURI(conf.Config.Mongo.Address))
	//if err!=nil{
	//	logger.Error("mongodb conn:%s fail,err：%v",conf.Config.Mongo.Address)
	//	panic(err)
	//}
	//err = p.mongoClient.Ping(ctx, readpref.Primary())
	//if err!=nil{
	//	logger.Error("mongodb ping :%s fail!",conf.Config.Mongo.Address)
	//}else{
	//	logger.Info("mongodb ping :%s success!",conf.Config.Mongo.Address)
	//}
	//
	//return p;
}

//卸载redis
func UnitRedis() {
	if redisClient != nil {
		//关闭redis连接
		if err := redisClient.Close(); err != nil {
			log.Logger().Error("redis close error:%v", err)
		} else {
			log.Logger().Info("redis close success!")
		}
	} else {
		log.Logger().Info("redis is nil,no need close !")
	}

	////close mongo
	//if p.mongoClient !=nil {
	//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//
	//	if err :=p.mongoClient.Disconnect(ctx);err !=nil{
	//		logger.Error("mongodb close error:%v",err)
	//	}else{
	//		logger.Info("mongodb close success!")
	//	}
	//}
}
