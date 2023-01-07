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
	"gorm.io/gorm"
)

var DEFAULT = "default"

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

func (p *InfraApp) WriteDB() *gorm.DB {
	if val, ok := p.MySqlMap[DEFAULT];ok{
		return val.Master
	}
	return nil
}
func (p *InfraApp) ReadDB() *gorm.DB {
	if val, ok := p.MySqlMap[DEFAULT];ok{
		if val.Slave != nil{
			return val.Slave
		}
		return val.Master
	}
	return nil
}

func (p *InfraApp) WriteDByName(name string) *gorm.DB {
	if val, ok := p.MySqlMap[name];ok{
		return val.Master
	}
	return nil
}
func (p *InfraApp) ReadDByName(name string) *gorm.DB {
	if val, ok := p.MySqlMap[name];ok{
		if val.Slave != nil{
			return val.Slave
		}
		return val.Master
	}
	return nil
}

func (p *InfraApp) Mongo() *mongo.Client{
	if val, ok := p.MongoMap[DEFAULT];ok{
		return val
	}
	return nil
}

func (p *InfraApp) MongoByName(name string) *mongo.Client{
	if val, ok := p.MongoMap[name];ok{
		return val
	}
	return nil
}

func (p *InfraApp) Redis() *redis.Client  {
	if val, ok := p.RedisMap[DEFAULT];ok{
		return val
	}
	return nil
}

func (p *InfraApp) RedisByName(name string) *redis.Client  {
	if val, ok := p.RedisMap[name];ok{
		return val
	}
	return nil
}