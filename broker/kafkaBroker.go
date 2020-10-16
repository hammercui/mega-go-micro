/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2020/9/11
 * Time: 20:37
 * Mail: hammercui@163.com
 *
 */
package broker

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/kafka/v2"
	"time"
	"wfServerMicro/infra/conf"
	"wfServerMicro/infra/log"
)

func NewKafkaBroker() broker.Broker {
	appConf := conf.GetConf().AppConf
	if(appConf.Env == conf.AppEnv_local){
		return nil
	}
	kafkaConf := conf.GetConf().KafkaConf
	//初始化broker
	//v1
	//bro := kafka.NewBroker(func(o *broker.Options) {
	//	o.Addrs = kafkaConf.Addrs
	//})
	//v2
	sConf := kafka.DefaultBrokerConfig
	//init连接超时时间2s
	sConf.Net.DialTimeout = 2 * time.Second
	bro := kafka.NewBroker(kafka.BrokerConfig(sConf))

	if err := bro.Init(func(o *broker.Options) {
		o.Addrs = kafkaConf.Addrs
	}); err != nil {
		log.Logger().Errorf("kafka broker启动失败,conf:%v,err:%+v", kafkaConf, err)
		return nil
	}
	if err := bro.Connect(); err != nil {
		log.Logger().Errorf("kafka broker连接失败,conf:%v,err:%+v", kafkaConf, err)
		return nil
	} else {
		log.Logger().Info("kafka broker连接成功,conf:%v", kafkaConf)
		return bro
	}
}
