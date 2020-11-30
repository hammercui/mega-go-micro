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
	"github.com/hammercui/mega-go-micro/log"
	"time"
)

func main() {
	app := infra.InitApp()

	//log.Logger().Info("22222222222222222222222")
	//启动http

	log.Logger().Info(app.Broker.String())
	time.Sleep(10000)
	//启动rpc
	// rpc.Start(app)
}