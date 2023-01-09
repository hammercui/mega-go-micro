/*
@Desc : 2020/8/24 9:47
@Version : 1.0.0
@Time : 2020/8/24 9:47
@Author : hammercui
@File : BaseService
@Company: Sdbean
*/
package infra

type BaseService struct {
	App *InfraApp
	id int
}

func NewBaseService(app *InfraApp) *BaseService {
	return &BaseService{App: app}
}

type AtyAutoConf struct {
	Desc string `json:"desc"`
	Val  string `json:"val"`
}