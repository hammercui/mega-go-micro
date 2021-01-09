/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/1/8
 * Time: 10:49
 * Mail: hammercui@163.com
 *
 */
package handler

import (
	"context"
	infra "github.com/hammercui/mega-go-micro"
	pbGo "github.com/hammercui/mega-go-micro/demo/proto/pbGo"
)

type DemoService struct {
	base *infra.BaseService
}

func (d DemoService) Info(ctx context.Context, req *pbGo.CommReq, resp *pbGo.DemoResp) error {
	resp.Code = 1
	resp.Msg = "ok"
	return nil
}

func NewDemoService(app *infra.InfraApp) *DemoService {
	return &DemoService{
		base: infra.NewBaseService(app),
	}
}
