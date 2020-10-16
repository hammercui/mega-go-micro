/*
@Desc : 日志文件行号
@Version : 1.0.0
@Time : 2020/8/25 15:16
@Author : hammercui
@File : logLineHook
@Company: Sdbean
*/
package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type LineHook struct {
	Field     string
	Skip      int
	levels    []logrus.Level
	Formatter func(file, function string, line int) string
}
func (hook *LineHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *LineHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = hook.Formatter(findCaller(hook.Skip))
	return nil
}

func NewLineHook() *LineHook {
	//appConfig := conf.GetConf().AppConf
	levelArray := []logrus.Level{
		logrus.InfoLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
		logrus.DebugLevel, logrus.WarnLevel,
	}
	//非prod环境启动debug和warn
	//if appConfig.Env != conf.AppEnv_prod {
	//	levelArray = append(levelArray, logrus.DebugLevel, logrus.WarnLevel, )
	//}
	hook := LineHook{
		Field:  "source",
		Skip:   5,
		levels: levelArray,
		Formatter: func(file, function string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		},
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	hook.Field = "line"

	return &hook
}

func findCaller(skip int) (string, string, int) {
	var (
		pc       uintptr
		file     string
		function string
		line     int
	)
	for i := 0; i < 10; i++ {
		pc, file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return file, function, line
}

func getCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return pc, file, line
}