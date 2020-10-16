package infra

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
	"github.com/hammercui/mega-go-micro/log"
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

func (p *BaseDao) SelectAll(outs interface{}, sqlStr string, values ...interface{}) error {
	type1 := reflect.TypeOf(outs)
	if type1.Kind() != reflect.Ptr {
		log.Logger().Error("第一个参数必须是指针,sql:",sqlStr)
		return errors.New("第一个参数必须是指针")
	}
	type2 := type1.Elem()	// 解指针后的类型
	if type2.Kind() != reflect.Slice {
		log.Logger().Error("第一个参数必须指向切片,sql:",sqlStr)
		return errors.New("第一个参数必须指向切片")
	}
	type3 := type2.Elem()
	if type3.Kind() != reflect.Ptr {
		log.Logger().Error("切片元素必须是指针类型,sql:",sqlStr)
		return errors.New("切片元素必须是指针类型")
	}

	rows, err := app.ReadOnlyDB.Raw(sqlStr, values...).Rows()
	defer rows.Close()
	if err != nil {
		log.Logger().Error("sql err:", err)
	}

	for rows.Next() {
		//  type3.Elem()是User, elem是*User
		elem := reflect.New(type3.Elem()) //type1解指针 相当于User,此时新建了User
		// 传入*User
		err := app.ReadOnlyDB.ScanRows(rows, elem.Interface())
		if err != nil {
			log.Logger().Error("gorm err:", err)
			continue
		}
		// reflect.ValueOf(result).Elem()是[]*User，Elem是*User，newSlice是[]*User
		newSlice := reflect.Append(reflect.ValueOf(outs).Elem(), elem)
		// 扩容后的slice赋值给*outs
		reflect.ValueOf(outs).Elem().Set(newSlice)
	}
	return nil
}
