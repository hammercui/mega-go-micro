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
	"github.com/hammercui/mega-go-micro/v2/base"
	pbGo "github.com/hammercui/mega-go-micro/v2/demo/proto/pbGo"
	"github.com/hammercui/mega-go-micro/v2/tool"
	"testing"
	"time"
)

func TestStartHttp(t *testing.T) {
	app := base.InitApp()
	fmt.Println(app.HttpRunning)
}

func Test_tool_PostJson(t *testing.T) {
	app := base.InitApp()
	fmt.Println(app.HttpRunning)
	var req = pbGo.CommReq{
		UserNo: "100",

	}
	var resp pbGo.DemoResp
	if err := tool.PostJson("http://localhost:8858/demo/info", &req, &resp); err == nil {
		t.Log(&resp)
	} else {
		t.Fatal(err)
	}
	time.Sleep(2*time.Second)
}
