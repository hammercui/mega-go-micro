/*
@Desc : 2020/8/24 9:47
@Version : 1.0.0
@Time : 2020/8/24 9:47
@Author : hammercui
@File : BaseService
@Company: Sdbean
*/
package infra

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/hammercui/mega-go-micro/log"
	"reflect"
)

type BaseService struct {
	App *InfraApp
}

func NewBaseService(app *InfraApp) *BaseService {
	return &BaseService{App: app}
}

type AtyAutoConf struct {
	Desc string `json:"desc"`
	Val  string `json:"val"`
}

//统一查询redis hash
func (p *BaseService) GetHashByKey(hashName string, key string, retModal interface{}) error {
	redisValue, err := p.App.RedisClient.HGet(hashName, key).Result()
	if err != nil {
		log.Logger().Errorf("redis hget:%s ,key:%s,err:%+v", hashName, key, err)
		return err
	}
	//string 转 struct
	err = json.Unmarshal([]byte(redisValue), retModal)
	if err != nil {
		log.Logger().Errorf("redis hget:%s,decode json err:%v,redisValue:%s", hashName, err, redisValue)
		return err
	}
	//logger.Info("redis hget:%s,key:%s,value:%+v", hashName, key, retModal)
	return nil
}

func (p *BaseService) GetHashStrByKey(hashName string, key string) string {
	redisValue, err := p.App.RedisClient.HGet(hashName, key).Result()
	if err != nil && err.Error() != redis.Nil.Error() {
		log.Logger().Errorf("redis hget:%s ,key:%s,err:%+v", hashName, key, err)
		return ""
	}
	return redisValue
}

//统一存储redis hash
func (p *BaseService) SetHashByKey(hashName string, key string, saveModal interface{}) error {
	str2 := saveModal
	t1 := reflect.TypeOf(saveModal)
	//入参不是string
	if t1.Kind() != reflect.String {
		b, err := json.Marshal(saveModal)
		if err != nil {
			log.Logger().Errorf("redis hset:%s,key:%s,value:%v,encode json err:%+v", hashName, key, saveModal, err)
			return err
		}
		str2 = string(b)
	}
	//logger.Debug("redis save hash:%s  key:%s,value:%s", hash,key, str2)
	_, err := p.App.RedisClient.HSet(hashName, key, str2).Result()
	if err != nil {
		log.Logger().Errorf("redis save hash:%s key:%s,value:%s  err:%+v", hashName, key, str2, err)
		return err
	}
	return nil
}

//统一删除redis hash
func (p *BaseService) DelHashByKey(hashName string, key string) error {
	_, err := p.App.RedisClient.HDel(hashName, key).Result()
	if err != nil {
		log.Logger().Errorf("redis del hash:%s key:%s,err:%+v", hashName, key, err)
		return err
	} else {
		log.Logger().Infof("redis del hash:%s key:%s,success!", hashName, key)
	}
	return nil
}

//统一查询redis
func (p *BaseService) GetByKey(key string, retModal interface{}) error {
	redisValue, err := p.App.RedisClient.Get(key).Result()
	if err != nil {
		log.Logger().Errorf("redis get key:%s,err:%+v", key, err)
		return err
	}
	//string 转 struct
	err = json.Unmarshal([]byte(redisValue), retModal)
	if err != nil {
		log.Logger().Errorf("redis get:%s,decode json err:%v,redisValue:%s", key, err, redisValue)
		return err
	}
	//logger.Info("redis hget:%s,key:%s,value:%+v", hashName, key, retModal)
	return nil
}

func (p *BaseService) GetStringByKey(key string) string {
	redisValue, err := p.App.RedisClient.Get(key).Result()
	if err != nil && err.Error() != redis.Nil.Error() {
		log.Logger().Errorf("redis get key:%s,err:%+v", key, err)
		return ""
	}
	return redisValue
}

//统一存储redis
func (p *BaseService) SetByKey(key string, saveModal interface{}) error {
	str2 := saveModal
	t1 := reflect.TypeOf(saveModal)
	//入参不是string
	if t1.Kind() != reflect.String {
		b, err := json.Marshal(saveModal)
		if err != nil {
			log.Logger().Errorf("redis set key:%s,value:%v,encode json err:%+v", key, saveModal, err)
			return err
		}
		str2 = string(b)
	}
	//logger.Debug("redis save hash:%s  key:%s,value:%s", hash,key, str2)
	_, err := p.App.RedisClient.Set(key, str2, 0).Result()
	if err != nil {
		log.Logger().Errorf("redis save key:%s,value:%s  err:%+v", key, str2, err)
		return err
	}
	return nil
}

//统一存储redis
func (p *BaseService) SetStringByKey(key string, saveModal string) error {
	//logger.Debug("redis save hash:%s  key:%s,value:%s", hash,key, str2)
	_, err := p.App.RedisClient.Set(key, saveModal, 0).Result()
	if err != nil {
		log.Logger().Errorf("redis save key:%s,value:%s  err:%+v", key, saveModal, err)
		return err
	}
	return nil
}

//统一从队尾获得
func (p *BaseService) PopListR(key string, retModal interface{}) error {
	redisValue, err := p.App.RedisClient.RPop(key).Result()

	if err != nil {
		//logger.Error("redis rpop key:%s ,err:%v", key, err)
		return err
	}
	if len(redisValue) <= 0 {
		retModal = nil
		return nil
	}

	//string 转 struct
	err = json.Unmarshal([]byte(redisValue), retModal)
	if err != nil {
		log.Logger().Errorf("redis rpop:%s,decode json err:%v,redisValue:%s", key, err, redisValue)
		return err
	}
	//logger.Info("redis hget:%s,key:%s,value:%+v", hashName, key, retModal)
	return nil
}

//统一插入队头
func (p *BaseService) PushListL(key string, inModel interface{}) error {
	b, err := json.Marshal(inModel)
	if err != nil {
		log.Logger().Errorf("redis lpush key:%s,value:%v,encode json err:%+v", key, inModel, err)
		return err
	}
	str2 := string(b)
	_, err = p.App.RedisClient.LPush(key, str2).Result()
	if err != nil {
		log.Logger().Errorf("redis lpush key:%s,value:%s  err:%+v", key, str2, err)
		return err
	}
	return nil
}

func (p *BaseService) GetRedisClientByName(name string) *redis.Client {
	return p.App.GetRedisClient(name)
}
