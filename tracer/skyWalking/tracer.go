/**
 * Description:链路追踪sky-walking
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/2/27
 * Time: 16:24
 * Mail: hammercui@163.com
 *
 */
package skyWalking

import (
	"fmt"
	"github.com/hammercui/go2sky"
	"github.com/hammercui/go2sky/reporter"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
)

//新建链路追踪实例
func NewSkyTracer() *go2sky.Tracer {
	appConf := conf.GetConf().AppConf
	if appConf.Env == conf.AppEnv_local {
		return nil
	}
	webName := fmt.Sprintf("%s-%s-api-%s", appConf.Group, appConf.Name, appConf.Env)
	webId := fmt.Sprintf("%s-%s", webName, appConf.NodeId)
	skyWalkingAddr := "172.25.220.245:11800"
	if addr, ok := appConf.Custom["skyAddr"]; ok && len(addr) > 0 {
		skyWalkingAddr = addr
	} else {
		//配置文件没有skyAddr,不建立连接
		return nil
	}
	r, err := reporter.NewGRPCReporter(skyWalkingAddr)
	if err != nil {
		log.Logger().Errorf("skyWalking create reporter error:%V", err)
		return nil
	}
	tracer, err2 := go2sky.NewTracer(webName, go2sky.WithReporter(r), go2sky.WithInstance(webId))
	if err2 != nil {
		log.Logger().Errorf("skyWalking create tracer error:%V", err2)
	}

	log.Logger().Infof("skyWalking init success!  addr: %s", skyWalkingAddr)
	return tracer
}
