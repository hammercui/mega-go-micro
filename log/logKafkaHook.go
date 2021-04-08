/*
@Desc : 返回kafka hook
@Version : 1.0.0
@Time : 2020/8/24 15:11
@Author : hammercui
@File : logKafkaHook
@Company: Sdbean
*/
package log

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/sirupsen/logrus"
	"time"
)

type KafkaHook struct {
	// Id of the hook
	id string

	// Log levels allowed
	levels []logrus.Level

	// Log entry formatter
	formatter logrus.Formatter

	// sarama.AsyncProducer
	producer sarama.AsyncProducer
}

// Create a new KafkaHook.
func NewKafkaHook(id string, levels []logrus.Level, formatter logrus.Formatter, brokers []string) (*KafkaHook, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokers, kafkaConfig)

	if err != nil {
		return nil, err
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			Logger().Errorf("Failed to send log entry to kafka: %v\n", err)
		}
	}()

	hook := &KafkaHook{
		id,
		levels,
		formatter,
		producer,
	}

	return hook, nil
}

func (hook *KafkaHook) Id() string {
	return hook.id
}

func (hook *KafkaHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *KafkaHook) Fire(entry *logrus.Entry) error {
	// Check time for partition key
	var partitionKey sarama.ByteEncoder

	// Get field time
	t, _ := entry.Data["time"].(time.Time)

	// Convert it to bytes
	b, err := t.MarshalBinary()

	if err != nil {
		return err
	}

	partitionKey = sarama.ByteEncoder(b)

	// Check topics
	var topics []string

	if ts, ok := entry.Data["topics"]; ok {
		//fmt.Println("topics",ts)
		if topics, ok = ts.([]string); !ok {
			return errors.New("Field topics must be []string")
		}
	} else {
		return errors.New("Field topics not found")
	}

	// Format before writing
	b, err = hook.formatter.Format(entry)

	if err != nil {
		return err
	}

	value := sarama.ByteEncoder(b)

	//写入kafka调整为异步
	go func() {
		for _, topic := range topics {
			hook.producer.Input() <- &sarama.ProducerMessage{
				Key:   partitionKey,
				Topic: topic,
				Value: value,
			}
		}
	}()

	return nil
}

//获得kafka hook实例
func getKafkaHook() *KafkaHook {
	appConfig := conf.GetConf().AppConf
	//全部日志
	levelArray := []logrus.Level{
		logrus.InfoLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
		logrus.WarnLevel,
	}
	//非prod环境打印debug
	if appConfig.Env != conf.AppEnv_prod {
		levelArray = append(levelArray, logrus.DebugLevel)
	}
	hook, err := NewKafkaHook(
		"kh",
		levelArray,
		&logrus.JSONFormatter{},
		appConfig.KafkaHookAddrs,
	)
	if err != nil {
		Logger().Errorf("kafka err:%+v", err)
	}
	return hook
}
