/*
@Desc : 基础设施app对象
@Version : 1.0.0
@Time : 2020/8/22 16:02
@Author : hammercui
@File : infraApp
@Company: Sdbean
*/
package infra

import (
	"github.com/go-redis/redis"
	"github.com/hammercui/go2sky"
	"github.com/hammercui/mega-go-micro/mysql"
	"github.com/hammercui/mega-go-micro/watch"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"go.mongodb.org/mongo-driver/mongo"
)

type InfraApp struct {
	HttpRunning bool
	RpcRunning  bool
	RedisMap    map[string]*redis.Client //redis集合
	Reg         registry.Registry        //服务注册
	Selector    selector.Selector        //服务发现
	Broker      broker.Broker            //消息订阅与发布
	ConfWatch   *watch.ConfWatch         //配置中心
	//todo mongo连接

	Tracer   *go2sky.Tracer           //链路追踪
	MySqlMap map[string]*mysql.Client //mysql集合
	MongoMap map[string]*mongo.Client //mongo集合
	//redis客户端map
	//redisPoolMap map[string]*infraRedis.RedisPool
	//mysql只读连接池
	//readOnlyDBPoolMap map[string]*mysql.DBPool
	//mysql读写连接池
	//readWriteDBPoolMap map[string]*mysql.DBPool
}

//instance
var app *InfraApp
//
////指定名称的db链接入池
//func (p *InfraApp) SetReadOnlyDBPool(key string, addr string, dbConn *gorm.DB) {
//	if pool, ok := p.readOnlyDBPoolMap[key]; ok {
//		pool.PushDB(dbConn)
//	} else {
//		readDBPool := mysql.NewDBPoll(addr, dbConn)
//		p.readOnlyDBPoolMap[key] = readDBPool
//	}
//}
//
//func (p *InfraApp) GetReadOnlyDB(key string) *gorm.DB {
//	return p.readOnlyDBPoolMap[key].GetDB()
//}
//
//func (p *InfraApp) SetReadWriteDBPool(key string, addr string, dbConn *gorm.DB) {
//	if pool, ok := p.readWriteDBPoolMap[key]; ok {
//		pool.PushDB(dbConn)
//	} else {
//		readDBPool := mysql.NewDBPoll(addr, dbConn)
//		p.readWriteDBPoolMap[key] = readDBPool
//	}
//}
//
//func (p *InfraApp) GetReadWriteDB(key string) *gorm.DB {
//	return p.readWriteDBPoolMap[key].GetDB()
//}
//
//func (p *InfraApp) SetRedisPool(key string, addr string, dbIndex int, client *redis.Client) {
//	if pool, ok := p.redisPoolMap[key]; ok {
//		pool.PushClient(client)
//	} else {
//		redisPool := infraRedis.NewRedisPool(key, addr, dbIndex, client)
//		p.redisPoolMap[key] = redisPool
//	}
//}
//
//func (p *InfraApp) GetRedisClient(key string) *redis.Client {
//	return p.redisPoolMap[key].GetClient()
//}
