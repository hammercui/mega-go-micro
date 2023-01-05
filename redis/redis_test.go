package redis

import (
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
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
