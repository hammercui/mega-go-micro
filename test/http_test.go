/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2020/10/16
 * Time: 15:52
 * Mail: hammercui@163.com
 *
 */
package test

import (
	"fmt"
	infra "github.com/hammercui/mega-go-micro"
	"testing"
)

func TestStartHttp(t *testing.T) {
	app := infra.InitApp()
	fmt.Println(app.HttpRunning)
}