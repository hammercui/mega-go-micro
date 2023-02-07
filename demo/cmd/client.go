/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2020/12/19
 * Time: 14:28
 * Mail: hammercui@163.com
 *
 */
package main

import (
	"fmt"
	"github.com/hammercui/mega-go-micro/v2/base"
	"github.com/micro/go-micro/v2/config"
)

//如何通过selector发现并选择服务
func main() {
	app := base.InitApp()
	next, err := app.Selector.Select("test")
	if err != nil {
		return
	}
	node, err := next()
	if err != nil || node != nil {
		return
	}
	fmt.Println("可用服务地址", node.Address)

	//配置中心
	config.NewConfig()
}
