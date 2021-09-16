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
	infra "github.com/hammercui/mega-go-micro"
	"github.com/hammercui/mega-go-micro/demo/http"
	"github.com/hammercui/mega-go-micro/demo/rpc"
)

func main() {
	app := infra.InitAppWithOpts(&infra.AppOpts{
		IsConfWatchOn:  true,
		IsBrokerOn:     true,
		IsRedisOn:      true,
		IsMongoOn:      true,
		IsSqlOn:        true,
		IsSkyWalkingOn: true,
	})
	// 启动http服务
	go http.Start(app)
	// 启动rpc
	rpc.Start(app)
}
