/*
@Desc : infra日志功能
@Version : 1.0.0
@Time : 2020/8/24 13:27
@Author : hammercui
@File : log
@Company: Sdbean
*/
package log

import (
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	lr "github.com/micro/go-plugins/logger/logrus/v2"
	"github.com/sirupsen/logrus"
	"os"
	"github.com/hammercui/mega-go-micro/conf"
)

var logrusSingle *logrus.Entry

func InitLog() {
	//配置文件
	appConfig := conf.GetConf().AppConf
	nodeId := appConfig.NodeId

	//是否使用通用topic
	topic := fmt.Sprintf("werewolf-web-activity-%s-log", appConfig.Env)
	if customTopic,ok := appConfig.Custom["kafkaHookTopic"];ok && customTopic != ""{
		topic = customTopic
	}

	//加载日志,使用logrus
	logrusEntry := logrus.WithFields(logrus.Fields{
		"name":   appConfig.FullAppName,
		"nodeId": nodeId,
		"topics": []string{topic},
	})
	logrusIns := logrusEntry.Logger
	logrusIns.SetFormatter(&TerminalTextFormatter{
		IsTerminal:      true,
		TimestampFormat: "2006/01/02 15:04:05.999",
	})
	logrusIns.SetOutput(os.Stdout)

	logrusIns.SetLevel(logrus.DebugLevel) //日志级别
	logrusIns.AddHook(NewLineHook())
	if appConfig.Env != conf.AppEnv_local {
		logrusIns.AddHook(getKafkaHook())
	}
	logrusIns.AddHook(getWriteAllFileHook())   //全部日志
	logrusIns.AddHook(getWriteErrorFileHook()) //错误日志
	logrusSingle = logrusEntry

	////系统级默认日志
	l := lr.NewLogger(lr.WithLogger(logrusIns)).Fields(map[string]interface{}{
		"name":   appConfig.FullAppName,
		"nodeId": nodeId,
		"topics": []string{topic},
	})
	logger.DefaultLogger = l
	//打印配置
	logrusSingle.Info("consul配置", conf.GetConf().ConsulConf)
}

//获得日志实例
func Logger() *logrus.Entry {
	return logrusSingle
}
