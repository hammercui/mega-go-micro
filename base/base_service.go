/*
@Desc : 2020/8/24 9:47
@Version : 1.0.0
@Time : 2020/8/24 9:47
@Author : hammercui
@File : BaseService
@Company: Sdbean
*/
package base

type BaseService struct {
	app *InfraApp
	id int
}

func NewBaseService() BaseService {
	return BaseService{app: App()}
}

type AtyAutoConf struct {
	Desc string `json:"desc"`
	Val  string `json:"val"`
}