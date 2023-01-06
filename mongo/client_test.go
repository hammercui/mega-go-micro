package mongo

import (
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"testing"
)

func Test_newMongoClient(t *testing.T) {
	log.Set(log.DefaultLogrus())
	newMongoClient(&conf.MongoConf{
		Addr:     "mongodb://localhost:27017/?maxPoolSize=500&minPoolSize=10",
		DbName:   "",
		Username: "",
		Password: "",
		Enable:   true,
	})
}
