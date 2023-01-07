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
package tracer

import (
	"fmt"
	"github.com/hammercui/go2sky"
	"github.com/hammercui/go2sky/reporter"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
)

func InitTracer()  *go2sky.Tracer{
	log.Logger().Infof("-------tracer init console-------")
	_conf := conf.GetConf()
	if _conf.Tracer == nil{
		log.Logger().Infof("trace not config")
		return nil
	}
	if !_conf.Tracer.Enable{
		log.Logger().Infof("trace disable")
		return nil
	}
	if _conf.Tracer.TracerType == "skyWalking"{
		return newSkyTracer(_conf)
	}
	return nil
}

//新建链路追踪实例
func newSkyTracer(c *conf.Config) *go2sky.Tracer {
	webName := fmt.Sprintf("%s-%s-%s", c.App.Group, c.App.Name, c.App.Env)
	webId := fmt.Sprintf("%s-%s", webName, c.App.NodeId)
	r, err := reporter.NewGRPCReporter(c.Tracer.Addr)
	if err != nil {
		log.Logger().Errorf("skyWalking create reporter error:%v", err)
		return nil
	}
	tracer, err2 := go2sky.NewTracer(webName, go2sky.WithReporter(r), go2sky.WithInstance(webId))
	if err2 != nil {
		log.Logger().Errorf("skyWalking create tracer error:%v", err2)
	}

	log.Logger().Infof("skyWalking init success!  addr: %s", c.Tracer.Addr)
	return tracer
}
