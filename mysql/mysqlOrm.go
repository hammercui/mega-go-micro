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
	dbConn := NewMysqlConn(DefaultMysqlReadOnlyDsn())
	//存入map
	return dbConn
}

//新建默认可读可写mysql
func DefaultMysqlReadWrite() *gorm.DB {
	return NewMysqlConn(DefaultMysqlReadWriteDsn())
}

func DefaultMysqlReadOnlyDsn() string  {
	mysqlConf := conf.GetConf().MysqlConf
	return GenMysqlDsn(mysqlConf,true)
}

func DefaultMysqlReadWriteDsn() string  {
	mysqlConf := conf.GetConf().MysqlConf
	return GenMysqlDsn(mysqlConf,false)
}

func GenMysqlDsn(mysqlConf *conf.MysqlConf,isReadOnly bool) string {
	addr := mysqlConf.Addr
	if isReadOnly{
		addr = mysqlConf.ReadAddr
	}
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConf.Username,
		mysqlConf.Password,
		addr,
		mysqlConf.DbName,
		mysqlConf.Charset,
	)
	return dsn
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