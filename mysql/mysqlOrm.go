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

//var readOnlyDB, readWriteDB *gorm.DB

//新建默认只读mysql
func DefaultMysqlReadOnly() *gorm.DB {
	dbConn := NewMysqlConn(DefaultMysqlDsn())
	//存入map
	return dbConn
}

func DefaultMysqlReadWrite() *gorm.DB {
	return NewMysqlConn(DefaultMysqlDsn())
}

func DefaultMysqlDsn() string  {
	mysqlConf := conf.GetConf().MysqlConf
	return GenMysqlDsn(mysqlConf)
}

func GenMysqlDsn(mysqlConf *conf.MysqlConf) string {
	addr := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConf.Username,
		mysqlConf.Password,
		mysqlConf.ReadAddr,
		mysqlConf.DbName,
		mysqlConf.Charset,
	)
	return addr
}

func NewMysqlConn(addr string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		Logger: NewGormLog(conf.GetConf()),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		log.Logger().Errorf("mysql :%s connect error!%s", addr, err)
		os.Exit(0)
	}
	log.Logger().Infof("mysql :%s connect success!", addr)
	return db
}