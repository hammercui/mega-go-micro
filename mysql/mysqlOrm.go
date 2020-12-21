package mysql

import (
	"fmt"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

var readOnlyDB, readWriteDB *gorm.DB

//初始化只读mysql
func InitMysqlReadOnly() *gorm.DB {
	mysqlConf := conf.GetConf().MysqlConf
	readAddr := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConf.Username,
		mysqlConf.Password,
		mysqlConf.ReadAddr,
		mysqlConf.DbName,
		mysqlConf.Charset,
	)
	db, err := gorm.Open(mysql.Open(readAddr), &gorm.Config{
		Logger: NewGormLog(conf.GetConf()),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		log.Logger().Errorf("mysql readonly:%s connect error!%v", readAddr, err)
		os.Exit(0)
	}
	log.Logger().Infof("mysql readonly:%s connect success!", readAddr)
	readOnlyDB = db
	return readOnlyDB
}

func InitMysqlReadWrite() *gorm.DB {
	mysqlConf := conf.GetConf().MysqlConf
	readwriteAddr := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConf.Username,
		mysqlConf.Password,
		mysqlConf.Addr,
		mysqlConf.DbName,
		mysqlConf.Charset,
	)
	db, err := gorm.Open(mysql.Open(readwriteAddr), &gorm.Config{
		Logger: NewGormLog(conf.GetConf()),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		log.Logger().Errorf("mysql readwrite:%s connect err!%v", readwriteAddr, err)
		os.Exit(0)
	}
	log.Logger().Infof("mysql readwrite:%s connect success!", readwriteAddr)
	readWriteDB = db
	return readWriteDB
}

func UnitMysql() {

}
