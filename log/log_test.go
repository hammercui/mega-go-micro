package log

import (
	"github.com/hammercui/mega-go-micro/conf"
	"testing"
)

func TestInitLog(t *testing.T) {
	conf.InitConfig()
	InitLog()
}

func TestDefault(t *testing.T) {
	_logrus := DefaultLogrus()
	Set(_logrus)
	Logger().Infof("default success")
}