/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2020/11/30
 * Time: 10:36
 * Mail: hammercui@163.com
 *
 */
package main

import (
	"github.com/hammercui/mega-go-micro/v2/base"
	"github.com/hammercui/mega-go-micro/v2/demo/http"
)

func main() {
	//var env string
	//flag.StringVar(&env,"env","","-env 后面的值")
	//flag.Parse()
	//fmt.Println("当前环境为：",env)

	app := base.InitApp()
	// 启动http服务
	http.Start(app)
	// 启动rpc
	//rpc.Start(app)
}
