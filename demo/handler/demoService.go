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
	infra "github.com/hammercui/mega-go-micro/v2"
	pbGo "github.com/hammercui/mega-go-micro/v2/demo/proto/pbGo"
)

type DemoService struct {
	infra.BaseService
}

func (d DemoService) Info(ctx context.Context, req *pbGo.CommReq, resp *pbGo.DemoResp) error {
	resp.Code = 1
	resp.Msg = "ok"
	return nil
}

func NewDemoService(app *infra.InfraApp) *DemoService {
	return &DemoService{ infra.BaseService{App:app}}
}
