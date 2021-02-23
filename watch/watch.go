/**
 * Description:监听并合并配置变化
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2020/12/19
 * Time: 18:07
 * Mail: hammercui@163.com
 *
 */
package watch

import (
	"errors"
	"fmt"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	microConfig "github.com/micro/go-micro/v2/config"
	microSource "github.com/micro/go-micro/v2/config/source"
	consulSrc "github.com/micro/go-plugins/config/source/consul/v2"
	"os"
	"reflect"
)

//配置监听
type ConfWatch struct {
	consulSource microSource.Source
	config       microConfig.Config
	preKey       string
	env          string
}

//新建配置中心监听
func NewConfWatch() *ConfWatch {
	env := conf.GetConf().AppConf.Env
	key := "werewolf_conf"
	if len(conf.GetConf().ConfigCenter.ConfKey) != 0 {
		key = conf.GetConf().ConfigCenter.ConfKey
	}
	consulConfPre := fmt.Sprintf("/%s/%s", env, key)

	confWatch := &ConfWatch{
		preKey: key,
		env:    string(env),
	}

	//1 配置consul源
	confWatch.consulSource = consulSrc.NewSource(
		// optionally specify consul address; default to localhost:8500
		consulSrc.WithAddress(conf.GetConf().ConfigCenter.ConsulAddrs[0]),
		// optionally specify prefix; defaults to /micro/config
		consulSrc.WithPrefix(consulConfPre),
		//// optionally strip the provided prefix from the keys, defaults to false
		//consulSrc.StripPrefix(true),
	)
	//2 load源
	ins, err := microConfig.NewConfig()
	if err != nil {
		fmt.Println("config init error", err)
	}
	confWatch.config = ins
	if err := confWatch.config.Load(confWatch.consulSource); err != nil {
		log.Logger().Error("init config center error !", err)
		os.Exit(0)
	}
	log.Logger().Info("init config center success ! pre: ", consulConfPre)

	//3 merge配置
	confWatch.mergeConf()

	return confWatch
}

func (p *ConfWatch) mergeConf() {
	//mysql
	mysqlMap := make(map[string]string)
	if err := p.Get("mysql", &mysqlMap); err != nil {
		os.Exit(0)
	}
	log.Logger().Debug("mysql from config center: ", mysqlMap)
	if len(mysqlMap) != 0 {
		conf.GetConf().MysqlConf.Addr = fmt.Sprintf("%s:%s", mysqlMap["host"], mysqlMap["port"])
	}

	//readMysql
	readMysqlMap := make(map[string]string)
	if err := p.Get("readMysql", &readMysqlMap); err != nil {
		os.Exit(0)
	}
	log.Logger().Debug("readMysql from config center: ", readMysqlMap)
	if len(readMysqlMap) != 0 {
		conf.GetConf().MysqlConf.ReadAddr = fmt.Sprintf("%s:%s", mysqlMap["host"], mysqlMap["port"])
	}

	//redis
	var redisMap = []map[string]interface{}{}
	err := p.Get("redis", &redisMap)
	if err != nil {
		os.Exit(0)
	}
	log.Logger().Debug("redis from config center: ", redisMap)
	if len(redisMap) > 0 {
		var redisAddrs []string
		for _, item := range redisMap {
			redisAdds := fmt.Sprintf("%s:%v", item["host"], item["port"])
			redisAddrs = append(redisAddrs, redisAdds)
		}
		conf.GetConf().RedisConf.Sentinels = redisAddrs
	}

	// kafka
	var kafkaMap = []map[string]interface{}{}
	if err = p.Get("kafka", &kafkaMap); err == nil && len(kafkaMap) > 0 {
		log.Logger().Debug("kafka from config center: ", kafkaMap)
		var kafkaAddrs []string
		for _, item := range kafkaMap {
			kafkaAddr := fmt.Sprintf("%s:%v", item["host"], item["port"])
			kafkaAddrs = append(kafkaAddrs, kafkaAddr)
		}
		conf.GetConf().KafkaConf.Addrs = kafkaAddrs
	}

	//mongo
	var mongoMap = []map[string]interface{}{}
	if err = p.Get("mongodb", &mongoMap); err == nil && len(mongoMap) > 0 {
		log.Logger().Debug("mongodb from config center: ", mongoMap)
		var mongoAddrs []string
		addrUrl := "mongodb://"
		for _, item := range kafkaMap {
			addr := fmt.Sprintf("%s:%v", item["host"], item["port"])
			mongoAddrs = append(mongoAddrs, addr)
			addrUrl = fmt.Sprintf("%s%s,", addrUrl, addr)
		}
		conf.GetConf().MongoConf.Addr = addrUrl
	}
}

//从配置中心获得配置
func (p *ConfWatch) Get(key string, out interface{}) error {
	outType := reflect.TypeOf(out)
	if outType.Kind() != reflect.Ptr {
		log.Logger().Error("config center get out参数必须是指针")
		return errors.New("config center get out参数必须是指针")
	}
	outElemType := outType.Elem() // 解指针后的类型

	switch outElemType {
	//map[string]string
	case reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf("")):
		getMap := p.config.Get(p.env, p.preKey, key).StringMap(map[string]string{})
		getMapValue := reflect.ValueOf(getMap)
		keys := getMapValue.MapKeys()
		for _, k := range keys {
			val := getMapValue.MapIndex(k)
			//获得out指针指向的值
			reflect.ValueOf(out).Elem().SetMapIndex(k, val)
		}
	default:
		err := p.config.Get(p.env, p.preKey, key).Scan(out)
		if err != nil {
			log.Logger().Errorf("config center get key:%s err:%+v", key, err)
			return err
		}
	}
	return nil
}

type WatchCallback func(outConf interface{}, err error)

//监听配置变更
func (p *ConfWatch) Watch(key string, out interface{}, callback WatchCallback) error {
	outType := reflect.TypeOf(out)
	if outType.Kind() != reflect.Ptr {
		log.Logger().Error("config center watch out参数必须是指针")
		return errors.New("config center watch out参数必须是指针")
	}
	outElemType := outType.Elem() // 解指针后的类型

	w, err := p.config.Watch(p.env, p.preKey, key)
	if err != nil {
		log.Logger().Errorf("config center watch key:%s err:%+v", key, err)
		return err
	}

	go func() {
		// wait for next value
		v, err := w.Next()
		if err != nil {
			// do something
			log.Logger().Errorf("config center watch key:%s err:%+v", key, err)
			callback(nil, err)
		} else {
			switch outElemType {
			//map[string]string
			case reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf("")):
				getMap := v.StringMap(map[string]string{})
				getMapValue := reflect.ValueOf(getMap)
				keys := getMapValue.MapKeys()
				for _, k := range keys {
					val := getMapValue.MapIndex(k)
					//获得out指针指向的值
					reflect.ValueOf(out).Elem().SetMapIndex(k, val)
				}
				callback(out, nil)
			//
			default:
				err := v.Scan(out)
				if err != nil {
					log.Logger().Errorf("config center watch key:%s err:%+v", key, err)
					callback(nil, err)
				} else {
					callback(out, nil)
				}
			}
		}

		//重新发起监听
		p.Watch(key, out, callback)
	}()

	log.Logger().Debugf("config center start watch key:%s ", key)
	return nil
}
