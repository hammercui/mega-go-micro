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
	"github.com/hammercui/mega-go-micro/v2/conf"
	"github.com/micro/go-micro/v2/logger"
	lr "github.com/micro/go-plugins/logger/logrus/v2"
	"github.com/sirupsen/logrus"
	"os"
)

var logrusSingle *logrus.Entry

func InitLog() {
	//配置文件
	_conf := conf.GetConf()
	nodeId := _conf.App.NodeId
	fmt.Println("-------log init console-------")
	//日志使用的topic
	topic := fmt.Sprintf("default-%s-log", _conf.App.Env)
	if _conf.Log != nil && _conf.Log.KafkaHookTopic != "" {
		topic = _conf.Log.KafkaHookTopic
	}
	//加载日志,使用logrus
	_logrusEntry := logrus.WithFields(logrus.Fields{
		"name":   _conf.App.FullAppName,
		"nodeId": nodeId,
		"topics": []string{topic},
	})

	_logrusIns := _logrusEntry.Logger
	_logrusIns.SetFormatter(&TerminalTextFormatter{
		IsTerminal:      true,
		TimestampFormat: "2006/01/02 15:04:05.999",
	})
	_logrusIns.SetOutput(os.Stdout)
	_logrusIns.SetLevel(formLogLevel(_conf.Log.Level)) //日志级别
	_logrusIns.AddHook(NewLineHook())
	//kafka hook
	if _conf.Log.KafkaHookEnable {
		_logrusIns.AddHook(getKafkaHook())
	}
	_logrusIns.AddHook(getWriteAllFileHook())   //全部日志
	_logrusIns.AddHook(getWriteErrorFileHook()) //错误日志

	//micro框架层日志替换为logrus
	l := lr.NewLogger(lr.WithLogger(_logrusIns)).Fields(map[string]interface{}{
		"name":   _conf.App.FullAppName,
		"nodeId": nodeId,
		"topics": []string{topic},
	})
	logger.DefaultLogger = l

	//打印配置
	logrusSingle = _logrusEntry
	logrusSingle.Infof("-------log init start-------")
	logrusSingle.Infof("env: %s", _conf.App.Env)
	logrusSingle.Infof("nodeId: %s", _conf.App.NodeId)
	logrusSingle.Infof("ip: %s", _conf.App.IP)
	logrusSingle.Infof("logout: %s", _conf.Log.LogoutPath)
	logrusSingle.Infof("init log success! %+v", _conf.Log)
}

func formLogLevel(level string) logrus.Level {
	switch level {
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	default:
		return logrus.TraceLevel
	}
}

// 获得日志实例
func Logger() *logrus.Entry {
	return logrusSingle
}

func Set(entry *logrus.Entry) {
	if logrusSingle == nil {
		logrusSingle = entry
	}
}

func DefaultLogrus() *logrus.Entry {
	_logrusEntry := logrus.WithFields(logrus.Fields{})
	_logrusIns := _logrusEntry.Logger
	l := lr.NewLogger(lr.WithLogger(_logrusIns)).Fields(map[string]interface{}{})
	logger.DefaultLogger = l
	return _logrusEntry
}
