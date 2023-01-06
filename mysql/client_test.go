package mysql

import (
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"testing"
)

func Test_newMysqlClient(t *testing.T){
	log.Set(log.DefaultLogrus())
	client := newMysqlClient("test",&conf.MysqlReadWriteConf{
		Master: &conf.MysqlConf{
			DSN:           "mega:mega@tcp(localhost:3306)/mega?charset=utf8mb4&parseTime=True&loc=Local",
			WarnThreshold: 2,
			DebugInfo:     true,
		},
		Slave:  &conf.MysqlConf{
			DSN:           "mega:mega@tcp(localhost:3306)/mega?charset=utf8mb4&parseTime=True&loc=Local",
			WarnThreshold: 2,
			DebugInfo:     true,
		},
		Enable: true,
	})
	if client.Name != "test" {
		t.Fatal("new client err")
	}
	if client.Master == nil && client.Slave == nil{
		t.Fatal("new client master slave all nil")
	}
}

func Test_newMysqlConn(t *testing.T) {
	log.Set(log.DefaultLogrus())
	newMysqlConn(&conf.MysqlConf{
		DSN:           "mega:mega@tcp(localhost:3306)/mega?charset=utf8mb4&parseTime=True&loc=Local",
		WarnThreshold: 2,
		DebugInfo:     true,
	})
}
