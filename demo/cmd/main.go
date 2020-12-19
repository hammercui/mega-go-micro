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
	"fmt"
	infra "github.com/hammercui/mega-go-micro"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	app := infra.InitApp()

	appConf := conf.GetConf().AppConf

	//启动http
	// Run the server
	service := web.NewService(
		web.Name("test"),
		web.Id("1"),
		web.Registry(app.Reg),
		web.Address(fmt.Sprintf("%s:%d", appConf.Ip, appConf.HttpPort)),
		web.Metadata(map[string]string{
			"version": "1.0.0",
			"tags":    "werewolf,web,activity,api",
		}),
	)
	// 启动http服务
	if err := service.Run(); err != nil {
		log.Logger().Error("gin start router err:", err)
	}

	// 启动rpc
	// rpc.Start(app)
}
