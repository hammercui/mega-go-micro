/*
@Desc : 2020/8/24 9:47
@Version : 1.0.0
@Time : 2020/8/24 9:47
@Author : hammercui
@File : BaseController
@Company: Sdbean
*/
package infra

//controller基础
type BaseController struct {
	App *InfraApp
}

func NewBaseController(app *InfraApp) *BaseController {
	return &BaseController{App: app}
}
