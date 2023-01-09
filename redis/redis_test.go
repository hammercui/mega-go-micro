package redis

import (
	"github.com/hammercui/mega-go-micro/v2/conf"
	"github.com/hammercui/mega-go-micro/v2/log"
	"testing"
)

func Test_InitRedis(t *testing.T) {
	conf.InitConfig()
	log.InitLog()
	redis := InitRedis()
	if len(redis) == 0 {
		t.Fatal("redis is empty")
	}
}
