package log

import (
	"github.com/hammercui/mega-go-micro/conf"
	"testing"
)

func TestInitLog(t *testing.T) {
	conf.InitConfig()
	InitLog()
}
