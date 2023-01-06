package mysql

import (
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/**
mysql connection client,include master slave
*/

type Client struct {
	Master *gorm.DB
	Slave  *gorm.DB
	Name   string
}

func InitMysql() map[string]*Client {
	_map := make(map[string]*Client)
	if conf.GetConf().MysqlMap == nil || len(conf.GetConf().MysqlMap) == 0 {
		return _map
	}

	log.Logger().Infof("-------mysql init console-------")
	for k, v := range conf.GetConf().MysqlMap {
		if v.Enable {
			_map[k] = newMysqlClient(k, v)
		}
	}
	return _map
}

func newMysqlClient(name string, c *conf.MysqlReadWriteConf) *Client {
	log.Logger().Infof("create mysql client: %s", name)
	_client := &Client{
		Master: nil,
		Slave:  nil,
		Name:   name,
	}
	_client.Master = newMysqlConn(c.Master)
	if c.Slave != nil {
		_client.Slave = newMysqlConn(c.Slave)
	}
	return _client
}

func newMysqlConn(c *conf.MysqlConf) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: NewGormLog(c),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		log.Logger().Errorf("mysql  connect error! dsn:%s,err:%v", c.DSN, err)
		panic(err)
	}
	log.Logger().Infof("mysql connect success! dsn:%s", c.DSN)
	return db
}
