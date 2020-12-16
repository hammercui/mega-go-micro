package infra

import (
	"database/sql"
	"errors"
	"github.com/hammercui/mega-go-micro/log"
	"gorm.io/gorm"
	"reflect"
)

type BaseDao struct {
	app *InfraApp
}

func NewBaseDao(app *InfraApp) *BaseDao {
	return &BaseDao{app: app}
}

//update delete 操作
func (p *BaseDao) Exec(sqlStr string, values ...interface{}) {
	app.ReadWriteDB.Exec(sqlStr, values...)
}

//
func (p *BaseDao) GetReadonlyDB() *gorm.DB {
	return app.ReadOnlyDB
}

//
func (p *BaseDao) GetReadWriteDB() *gorm.DB {
	return app.ReadWriteDB
}

//查询自定义返回对象
func (p *BaseDao) SelectCustom(out []interface{}, sqlStr string, values ...interface{}) error {
	type1 := reflect.TypeOf(out)
	if type1.Kind() != reflect.Slice {
		log.Logger().Error("第一个参数必须是interface切片")
		return errors.New("第一个参数必须是interface切片")
	}
	if len(out) == 0 {
		log.Logger().Error("第一个参数长度不能为空")
		return errors.New("第一个参数长度不能为空")
	}
	var row *gorm.DB
	if len(values) > 0 {
		row = app.ReadOnlyDB.Raw(sqlStr, values...)
	} else {
		row = app.ReadOnlyDB.Raw(sqlStr)
	}

	if row.Row() == nil {
		return nil
	}

	if row.Row().Err() != nil {
		log.Logger().Error("SelectCustom sql error:", row.Row().Err(), " | sqlStr: ", sqlStr)
		return row.Row().Err()
	}

	if err := row.Row().Scan(out...); err != nil && err != sql.ErrNoRows {
		log.Logger().Error("SelectCustom sql error:", err, " | sqlStr: ", sqlStr)
		return err
	}
	return nil
}

//查询struct对象
func (p *BaseDao) SelectOne(out interface{}, sqlStr string, values ...interface{}) error {
	type1 := reflect.TypeOf(out)
	if type1.Kind() != reflect.Ptr {
		log.Logger().Error("第一个参数必须是指针")
		return errors.New("第一个参数必须是指针")
	}
	var row *gorm.DB
	if len(values) > 0 {
		row = app.ReadOnlyDB.Raw(sqlStr, values...)
	} else {
		row = app.ReadOnlyDB.Raw(sqlStr)
	}
	row.Scan(out)
	return nil
}

//查询struct列表
func (p *BaseDao) SelectAll(outs interface{}, sqlStr string, values ...interface{}) error {
	type1 := reflect.TypeOf(outs)
	if type1.Kind() != reflect.Ptr {
		log.Logger().Error("第一个参数必须是指针,sql:", sqlStr)
		return errors.New("第一个参数必须是指针")
	}
	type2 := type1.Elem() // 解指针后的类型
	if type2.Kind() != reflect.Slice {
		log.Logger().Error("第一个参数必须指向切片,sql:", sqlStr)
		return errors.New("第一个参数必须指向切片")
	}
	type3 := type2.Elem()
	if type3.Kind() != reflect.Ptr {
		log.Logger().Error("切片元素必须是指针类型,sql:", sqlStr)
		return errors.New("切片元素必须是指针类型")
	}

	rows, err := app.ReadOnlyDB.Raw(sqlStr, values...).Rows()
	defer rows.Close()
	if err != nil {
		log.Logger().Error("SelectAll sql err:", err, " | sqlStr: ", sqlStr)
	}

	for rows.Next() {
		//  type3.Elem()是User, elem是*User
		elem := reflect.New(type3.Elem()) //type1解指针 相当于User,此时新建了User
		// 传入*User
		err := app.ReadOnlyDB.ScanRows(rows, elem.Interface())
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
