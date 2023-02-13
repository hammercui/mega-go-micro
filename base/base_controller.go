/*
@Desc : 2020/8/24 9:47
@Version : 1.0.0
@Time : 2020/8/24 9:47
@Author : hammercui
@File : BaseController
@Company: Sdbean
*/
package base

//controller基础
type BaseController struct {
	app *InfraApp
}

func NewBaseController() BaseController {
	return BaseController{app: App()}
}

