/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/1/8
 * Time: 10:08
 * Mail: hammercui@163.com
 *
 */
package http

import (
	"fmt"
	infra "github.com/hammercui/mega-go-micro"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/demo/handler"
	"github.com/hammercui/mega-go-micro/demo/http/middleware"
	pbGo "github.com/hammercui/mega-go-micro/demo/proto/pbGo"
	"github.com/hammercui/mega-go-micro/http/gin"
	"github.com/hammercui/mega-go-micro/log"
	"github.com/micro/go-micro/v2/web"
	"os"
	"strconv"
	"time"
)

func Start(app *infra.InfraApp) {
	ginServer := gin.NewMegaGinServer(app,
		middleware.Logger(),
		middleware.Recovery(),
	)

	//注册路由
	registerRouter(app, ginServer)

	appConf := conf.GetConf().App
	webName := fmt.Sprintf("%s-%s-api-%s", appConf.Group, appConf.Name, appConf.Env)
	webId := fmt.Sprintf("%s-%s", webName, appConf.NodeId)

	//注册服务发现
	httpService := web.NewService(
		web.Name(webName),
		web.Id(webId),
		web.Address(fmt.Sprintf("0.0.0.0:%d", appConf.HttpPort)),
		web.Handler(ginServer.Gin()),
		web.Registry(app.Reg),
		web.AfterStart(func() error {
			app.HttpRunning = true
			return nil
		}),
		web.BeforeStop(func() error {
			fmt.Sprintln("gin before stop")
			app.HttpRunning = false
			return nil
		}),
		web.AfterStart(func() error {
			fmt.Sprintln("gin after stop")
			return nil
		}),
		web.Version("1.0.0"),
		web.Metadata(map[string]string{
			"version": "1.0.0",
			"tags":    "werewolf,web,activity,api",
			"ip":      appConf.IP,
			"port":    strconv.Itoa(appConf.HttpPort),
		}), //元数据
		web.RegisterInterval(180*time.Second),
	)
	// 运行服务
	if err := httpService.Run(); err != nil {
		log.Logger().Error("gin start router err:", err)
	} else {
		os.Exit(0)
	}
}

func registerRouter(app *infra.InfraApp, ginServer *gin.GinServer) {
	//链路追踪
	if app.Tracer != nil {
		//ginRouter.Use(hammerHttp.SkyWalking(ginRouter, trace))
		ginServer.Gin().Use(gin.SkyWalkingMiddleware(ginServer.Gin(), app.Tracer))
	}

	//demo
	pbGo.RegisterDemoHandler(ginServer.Server(), handler.NewDemoService(app))
}
