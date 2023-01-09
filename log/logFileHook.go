/*
@Desc : 文件日志hook
@Version : 1.0.0
@Time : 2020/9/3 13:34
@Author : hammercui
@File : logFileHook
@Company: Sdbean
*/
package log

import (
	"fmt"
	"github.com/hammercui/mega-go-micro/v2/conf"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

func getWriteAllFileHook() *lfshook.LfsHook {
	_conf := conf.GetConf()
	logFileName := _conf.App.Name
	baseLogPath := path.Join(_conf.Log.LogoutPath, logFileName+".log")
	//全部writer
	writer, err := rotatelogs.New(
		path.Join(_conf.Log.LogoutPath, logFileName+".%Y-%m-%d"+".log"),
		rotatelogs.WithLinkName(baseLogPath),                                // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(_conf.Log.MaxDay*24)*time.Hour), // 文件最大保存时间 15天
		rotatelogs.WithRotationTime(24*time.Hour),                           // 日志切割时间间隔 1天
	)
	if err != nil {
		fmt.Printf("config local file system logger error. %v \n", errors.WithStack(err))
	}

	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}
	//非prod环境启动debug和warn
	if _conf.App.Env != conf.AppEnv_prod {
		writerMap[logrus.DebugLevel] = writer
		writerMap[logrus.WarnLevel] = writer
	}
	logOutHook := lfshook.NewHook(writerMap, &FileTextFormatter{
		TimestampFormat: "2006/01/02 15:04:05.999",
	})
	return logOutHook
}

func getWriteErrorFileHook() *lfshook.LfsHook {
	_conf := conf.GetConf()
	logFileName := _conf.App.Name
	baseErrLogPath := path.Join(_conf.Log.LogoutPath, logFileName+".err.log")

	errWriter, err := rotatelogs.New(
		path.Join(_conf.Log.LogoutPath, logFileName+".err.%Y-%m-%d"+".log"),
		rotatelogs.WithLinkName(baseErrLogPath),                             // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(_conf.Log.MaxDay*24)*time.Hour), // 文件最大保存时间 30天
		rotatelogs.WithRotationTime(24*time.Hour),                           // 日志切割时间间隔 1天
	)
	if err != nil {
		fmt.Printf("config local file system logger error. %v \n", errors.WithStack(err))
	}
	logOutHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.ErrorLevel: errWriter,
		logrus.FatalLevel: errWriter,
		logrus.PanicLevel: errWriter,
	}, &FileTextFormatter{
		TimestampFormat: "2006/01/02 15:04:05.999",
	})
	return logOutHook
}
