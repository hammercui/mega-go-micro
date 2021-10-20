package mysql

import (
	"fmt"
	"gorm.io/gorm"
)

type DBPool struct {
	addr  string
	pools []*gorm.DB
}

func  NewDBPoll(addr string,dbConn *gorm.DB) *DBPool {
	return &DBPool{
		addr: addr,
		pools: []*gorm.DB{
			dbConn,
		},
	}
}

func (p *DBPool) GetAddr() string  {
	return p.addr
}

func (p *DBPool) GetDB() *gorm.DB  {
	if len(p.pools) <1{
		panic(fmt.Sprintf("db connect is not exist,addr:%s",p.addr))
	}
	return p.pools[0]
}

func (p *DBPool) PushDB(dbConn *gorm.DB)  {
	p.pools = append(p.pools,dbConn)
}