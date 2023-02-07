package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hammercui/mega-go-micro/v2/log"
	"reflect"
)

type RedisService struct {
	BaseService
}

func NewRedisService() *RedisService {
	return &RedisService{
		BaseService{
			app: App(),
		},
	}
}


//统一查询redis hash
func (p *RedisService) HGetStrRedis(redisName string,hashName string, key string) (string,error) {
	client :=  p.app.RedisByName(redisName)
	if client == nil {
		return "",errors.New(fmt.Sprintf("redis[%s] is nil",redisName))
	}
	redisValue, err := client.HGet(hashName, key).Result()
	if err != nil && err.Error() != redis.Nil.Error() {
		log.Logger().Errorf("redis hget:%s ,key:%s,err:%+v", hashName, key, err)
		return "",err
	}
	return redisValue,nil
}
func (p *RedisService) HGetStr(hashName string, key string) (string,error){
	return p.HGetStrRedis(DEFAULT, hashName, key)
}
func (p *RedisService) HGetRedis(redisName string,hashName string, key string, retModal interface{}) error {
	redisValue, err := p.HGetStrRedis(redisName,hashName, key)
	if err != nil {
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
func (p *RedisService) HGet(hashName string, key string, retModal interface{}) error {
	return p.HGetRedis(DEFAULT,hashName,key,retModal)
}

//统一存储redis hash
func (p *RedisService) HSetRedis(redisName string,hashName string, key string, saveModal interface{}) error {
	client :=  p.app.RedisByName(redisName)
	if client == nil {
		return errors.New(fmt.Sprintf("redis[%s] is nil",redisName))
	}
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
	_, err := client.HSet(hashName, key, str2).Result()
	if err != nil {
		log.Logger().Errorf("redis save hash:%s key:%s,value:%s  err:%+v", hashName, key, str2, err)
		return err
	}
	return nil
}
func (p *RedisService) HSet(hashName string, key string, saveModal interface{}) error {
	return p.HSetRedis(DEFAULT,hashName,key,saveModal)
}

//统一删除redis hash
func (p *RedisService) HDelRedis(redisName string,hashName string, key string) error {
	client :=  p.app.RedisByName(redisName)
	if client == nil {
		return errors.New(fmt.Sprintf("redis[%s] is nil",redisName))
	}
	_, err := client.HDel(hashName, key).Result()
	if err != nil {
		log.Logger().Errorf("redis del hash:%s key:%s,err:%+v", hashName, key, err)
		return err
	} else {
		log.Logger().Infof("redis del hash:%s key:%s,success!", hashName, key)
	}
	return nil
}
func (p *RedisService) HDel(hashName string, key string) error {
	return p.HDelRedis(DEFAULT,hashName,key)
}

//统一查询redis
func (p *RedisService) Get(key string, retModal interface{}) error {
	return p.GetRedis(DEFAULT,key, retModal)
}
func (p *RedisService) GetRedis(redisName string,key string, retModal interface{}) error {
	redisValue, err := p.GetStrRedis(redisName,key)
	if err != nil {
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
func (p *RedisService) GetStr(key string) (string,error) {
	return p.GetStrRedis(DEFAULT,key)
}
func (p *RedisService) GetStrRedis(redisName string,key string) (string,error) {
	client :=  p.app.RedisByName(redisName)
	if client == nil {
		return "",errors.New(fmt.Sprintf("redis[%s] is nil",redisName))
	}
	redisValue, err := client.Get(key).Result()
	if err != nil && err.Error() != redis.Nil.Error() {
		log.Logger().Errorf("redis get key:%s,err:%+v", key, err)
		return "",err
	}
	return redisValue,nil
}

//统一存储redis
func (p *RedisService) Set(key string, saveModal interface{}) error {
	return p.SetRedis(DEFAULT,key,saveModal)
}
func (p *RedisService) SetRedis(redisName string,key string, saveModal interface{}) error {
	client :=  p.app.RedisByName(redisName)
	if client == nil {
		return errors.New(fmt.Sprintf("redis[%s] is nil",redisName))
	}
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
	_, err := client.Set(key, str2, 0).Result()
	if err != nil {
		log.Logger().Errorf("redis save key:%s,value:%s  err:%+v", key, str2, err)
		return err
	}
	return nil
}

//统一从队尾获得
func (p *RedisService) RPop(key string, retModal interface{}) error {
	return p.RPopRedis(DEFAULT,key,retModal)
}
func (p *RedisService) RPopRedis(redisName string,key string, retModal interface{}) error {
	client :=  p.app.RedisByName(redisName)
	if client == nil {
		return errors.New(fmt.Sprintf("redis[%s] is nil",redisName))
	}
	redisValue, err := client.RPop(key).Result()
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
func (p *RedisService) LPush(key string, inModel interface{}) error{
	return p.LPushRedis(DEFAULT,key,inModel)
}
func (p *RedisService) LPushRedis(redisName string,key string, inModel interface{}) error {
	client :=  p.app.RedisByName(redisName)
	if client == nil {
		return errors.New(fmt.Sprintf("redis[%s] is nil",redisName))
	}
	b, err := json.Marshal(inModel)
	if err != nil {
		log.Logger().Errorf("redis lpush key:%s,value:%v,encode json err:%+v", key, inModel, err)
		return err
	}
	str2 := string(b)
	_, err = client.LPush(key, str2).Result()
	if err != nil {
		log.Logger().Errorf("redis lpush key:%s,value:%s  err:%+v", key, str2, err)
		return err
	}
	return nil
}

