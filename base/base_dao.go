package base

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hammercui/mega-go-micro/v2/log"
	"gorm.io/gorm"
	"reflect"
)

/**
关系型数据库操作基类
*/
type BaseDao struct {
	app *InfraApp
}

func NewBaseDao() *BaseDao {
	return &BaseDao{app: App()}
}

//update delete by default db
func (p *BaseDao) Exec(sqlStr string, values ...interface{}) error {
	return p.ExecDB(DEFAULT,sqlStr,values...)
}
func (p *BaseDao) ExecDB(dbName string, sqlStr string, values ...interface{}) error   {
	db := p.app.WriteDByName(dbName)
	if db == nil{
		return errors.New(fmt.Sprintf("write db[%s] is nil",dbName))
	}
	db.Exec(sqlStr, values...)
	return nil
}

//插入orm对象
func (p *BaseDao) Insert(record interface{}) error {
	return p.InsertDB(DEFAULT,record)
}
func (p *BaseDao) InsertDB(dbName string,record interface{}) error {
	db := p.app.WriteDByName(dbName)
	if db == nil{
		return errors.New(fmt.Sprintf("write db[%s] is nil",dbName))
	}
	type1 := reflect.TypeOf(record)
	if type1.Kind() != reflect.Ptr {
		log.Logger().Errorf("record必须是指针,record:%+v", record)
		return errors.New("record必须是指针")
	}
	result := db.Create(record)
	if result.Error != nil {
		log.Logger().Errorf("insert record:%+v,error:%+v", record, result.Error)
		return result.Error
	}
	return nil
}

//查询自定义返回对象
func (p *BaseDao) SelectCustom(out []interface{}, sqlStr string, values ...interface{}) error   {
	return p.SelectCustomDB(DEFAULT,out,sqlStr,values...)
}
func (p *BaseDao) SelectCustomDB(dbName string, out []interface{}, sqlStr string, values ...interface{}) error {
	db := p.app.ReadDByName(dbName)
	if db == nil{
		return errors.New(fmt.Sprintf("read db[%s] is nil",dbName))
	}
	type1 := reflect.TypeOf(out)
	if type1.Kind() != reflect.Slice {
		log.Logger().Errorf("第一个参数必须是interface切片,sql: %s", sqlStr)
		return errors.New("第一个参数必须是interface切片")
	}
	if len(out) == 0 {
		log.Logger().Errorf("第一个参数长度不能为空,sql: %s", sqlStr)
		return errors.New("第一个参数长度不能为空")
	}
	var row *gorm.DB
	if len(values) > 0 {
		row = db.Raw(sqlStr, values...)
	} else {
		row = db.Raw(sqlStr)
	}

	if err := row.Row().Scan(out...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Logger().Errorf("SelectCustom sql error: %v, sqlStr: %s",err, sqlStr)
		return err
	}
	return nil
}

//查询struct对象
func (p *BaseDao) SelectOne(out interface{}, sqlStr string, values ...interface{}) error {
	return p.SelectOneDB(DEFAULT,out,sqlStr,values...)
}
func (p *BaseDao) SelectOneDB(dbName string,out interface{}, sqlStr string, values ...interface{}) error {
	db := p.app.ReadDByName(dbName)
	if db == nil{
		return errors.New(fmt.Sprintf("read db[%s] is nil",dbName))
	}
	type1 := reflect.TypeOf(out)
	if type1.Kind() != reflect.Ptr {
		log.Logger().Errorf("第一个参数必须是指针,sql: %s", sqlStr)
		return errors.New("第一个参数必须是指针")
	}
	var row *gorm.DB
	if len(values) > 0 {
		row = db.Raw(sqlStr, values...)
	} else {
		row = db.Raw(sqlStr)
	}
	row.Scan(out)
	return nil
}

//查询struct列表
func (p *BaseDao) SelectAll(outs interface{}, sqlStr string, values ...interface{}) error {
	return p.SelectAllDB(DEFAULT,outs,sqlStr,values...)
}
func (p *BaseDao) SelectAllDB(dbName string,outs interface{}, sqlStr string, values ...interface{}) error {
	db := p.app.ReadDByName(dbName)
	if db == nil{
		return errors.New(fmt.Sprintf("read db[%s] is nil",dbName))
	}
	type1 := reflect.TypeOf(outs)
	if type1.Kind() != reflect.Ptr {
		log.Logger().Errorf("第一个参数必须是指针,sql: %s", sqlStr)
		return errors.New("第一个参数必须是指针")
	}
	type2 := type1.Elem() // 解指针后的类型
	if type2.Kind() != reflect.Slice {
		log.Logger().Errorf("第一个参数必须指向切片,sql: %s", sqlStr)
		return errors.New("第一个参数必须指向切片")
	}
	type3 := type2.Elem()
	if type3.Kind() != reflect.Ptr {
		log.Logger().Errorf("切片元素必须是指针类型,sql: %s", sqlStr)
		return errors.New("切片元素必须是指针类型")
	}

	rows, err := db.Raw(sqlStr, values...).Rows()
	defer rows.Close()
	if err != nil {
		log.Logger().Error("SelectAll sql err:", err, " | sqlStr: ", sqlStr)
		return err
	}

	for rows.Next() {
		//  type3.Elem()是User, elem是*User
		elem := reflect.New(type3.Elem()) //type1解指针 相当于User,此时新建了User
		// 传入*User
		err := db.ScanRows(rows, elem.Interface())
		if err != nil {
			log.Logger().Error("SelectAll gorm err:", err, " | sqlStr: ", sqlStr)
			continue
		}
		// reflect.ValueOf(result).Elem()是[]*User，Elem是*User，newSlice是[]*User
		newSlice := reflect.Append(reflect.ValueOf(outs).Elem(), elem)
		// 扩容后的slice赋值给*outs
		reflect.ValueOf(outs).Elem().Set(newSlice)
	}
	return nil
}