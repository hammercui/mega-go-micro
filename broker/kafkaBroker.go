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
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/kafka/v2"
	"time"
)

func InitKafkaBroker() broker.Broker {
	_conf := conf.GetConf()
	if _conf.Kafka == nil {
		return nil
	}
	kafkaConf := _conf.Kafka
	if !kafkaConf.Enable{
		return nil
	}
	log.Logger().Info("-------kafka init console-------")
	return newKafka(kafkaConf)
}

func newKafka(c *conf.KafkaConf)  broker.Broker{
	sConf := kafka.DefaultBrokerConfig
	//init连接超时时间2s
	sConf.Net.DialTimeout = 2 * time.Second
	if c.DialTimeout > 0 {
		sConf.Net.DialTimeout = time.Duration(c.DialTimeout) * time.Second
	}
	bro := kafka.NewBroker(kafka.BrokerConfig(sConf))
	if err := bro.Init(func(o *broker.Options) {
		o.Addrs = c.Addrs
	}); err != nil {
		log.Logger().Errorf("kafka broker connect fail!,conf:%v,err:%+v", c, err)
		return nil
	}
	if err := bro.Connect(); err != nil {
		log.Logger().Errorf("kafka broker connect fail!,conf:%+v,err:%+v", c, err)
		return nil
	} else {
		log.Logger().Infof("kafka broker connect success!,addr:%s", c.Addrs)
		return bro
	}
}
