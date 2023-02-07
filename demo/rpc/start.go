/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/1/9
 * Time: 10:35
 * Mail: hammercui@163.com
 *
 */
package rpc

import (
	"fmt"
	"github.com/hammercui/mega-go-micro/v2/base"
	"github.com/hammercui/mega-go-micro/v2/conf"
	"github.com/hammercui/mega-go-micro/v2/demo/handler"
	pbGo "github.com/hammercui/mega-go-micro/v2/demo/proto/pbGo"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	"strconv"
)

func Start(app *base.InfraApp) {
	appConf := conf.GetConf().App
	rpcName := fmt.Sprintf("%s-%s-rpc-%s", appConf.Group, appConf.Name, appConf.Env)
	// New Service
	// 创建新的服务，这里可以传入其它选项。
	service := micro.NewService(
		micro.Registry(app.Reg),
		micro.Name(rpcName),
		micro.Address(fmt.Sprintf("%s:%d", "0.0.0.0", appConf.RpcPort)),
		micro.Metadata(map[string]string{
			"version": "1.0.0",
			"tags":    "werewolf,web,activity,rpc",
			"ip":      appConf.IP,
			"port":    strconv.Itoa(appConf.RpcPort),
		}),
		//micro.WrapHandler(skyWalking.NewHandlerWrapper(app.Tracer, "User-Agent")),
	)
	service.Init()

	// Register Handler
	registerHandler(service.Server(), app)

	// Run service
	// 运行服务
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func registerHandler(ser server.Server, app *base.InfraApp) {
	//demo
	pbGo.RegisterDemoHandler(ser, handler.NewDemoService(app))
}
