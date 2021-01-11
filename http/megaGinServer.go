/**
 * Description:封装的基于gin的http服务
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/1/7
 * Time: 17:39
 * Mail: hammercui@163.com
 *
 */
package http

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	infra "github.com/hammercui/mega-go-micro"
	"github.com/hammercui/mega-go-micro/log"
	"github.com/micro/go-micro/v2/server"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type GinServer struct {
	ginRouter *gin.Engine
	app       *infra.InfraApp
}

func (p *GinServer) NewSubscriber(s string, i interface{}, option ...server.SubscriberOption) server.Subscriber {
	panic("implement me")
}

func (p *GinServer) Subscribe(subscriber server.Subscriber) error {
	panic("implement me")
}

func (p *GinServer) Start() error {
	panic("implement me")
}

func (p *GinServer) Stop() error {
	panic("implement me")
}

func (p *GinServer) String() string {
	panic("implement me")
}

func (p *GinServer) Server() *GinServer {
	return p
}

func (p *GinServer) Gin() *gin.Engine {
	return p.ginRouter
}

//实例化megaGinServer
func NewMegaGinServer(app *infra.InfraApp, middlewares ...gin.HandlerFunc) *GinServer {
	gin.DisableConsoleColor()
	//gin设置模式
	gin.SetMode(gin.DebugMode)
	//初始化路由
	ginRouter := gin.New()
	//注册通用中间件
	for _, item := range middlewares {
		ginRouter.Use(item)
	}
	//健康检查
	ginRouter.GET("", healthResponse)
	ginRouter.GET("/actuator", healthResponse)
	ginRouter.GET("/actuator/health", healthResponse)

	//心跳
	ginRouter.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"code": 200,
			"sign": 200,
			"msg":  "pong!",
		})
	})

	return &GinServer{
		ginRouter: ginRouter,
		app:       app,
	}
}

func (a *GinServer) Init(option ...server.Option) error {
	panic("implement me")
}

func (a *GinServer) Options() server.Options {
	panic("implement me")
}

func (a *GinServer) Handle(handler server.Handler) error {
	return nil
}

func (a *GinServer) NewHandler(i interface{}, opts ...server.HandlerOption) server.Handler {
	//i为处理service
	//option包含
	options := server.HandlerOptions{
		Metadata: make(map[string]map[string]string),
	}

	for _, o := range opts {
		o(&options)
	}
	for k, v := range options.Metadata {
		log.Logger().Debugf("metadata k:%s,v:%v", k, v)
		//方法名
		funName := strings.Split(k, ".")[1]
		method := v["method"]
		path := v["path"]
		a.registerEndPointsV2Imp(funName, method, path, i)

	}
	return nil
}

//注册endpoint实现
func (p *GinServer) registerEndPointsV2Imp(funName string, method string, path string, service interface{}) {
	a := reflect.TypeOf(service)
	if m, ok := a.MethodByName(funName); ok {
		log.Logger().Infof("注册service路由成功(post):%s,类%s ,方法%s", path, a.String(), m.Name)
		a := func(c *gin.Context) {
			//形参数量
			parametersNum := m.Type.NumIn()
			parameters := make([]reflect.Value, parametersNum)
			parameters[0] = reflect.ValueOf(service)
			parameters[1] = reflect.ValueOf(c)

			//入参序列化
			reqPtr := m.Type.In(2)
			var req = reflect.New(reqPtr.Elem())
			if erro := bindJson(req.Interface(), c); erro != nil {
				return
			}
			parameters[2] = reflect.ValueOf(req.Interface())

			//出参
			respPtr := m.Type.In(3)
			var resp = reflect.New(respPtr.Elem())
			parameters[3] = reflect.ValueOf(resp.Interface())

			//执行函数
			result := m.Func.Call(parameters)

			//返回结果
			for _, item := range result {
				if item.Interface() == nil {
					c.JSON(http.StatusOK, resp.Interface())
				} else {
					dieFail(item.Interface().(error), c)
				}
			}
		}

		switch method {
		case "POST":
			p.ginRouter.POST(path, a)
		case "GET":
			p.ginRouter.GET(path, a)
		}

	} else {
		log.Logger().Errorf("注册service路由失败(post):%s,类%s ,未实现方法%s", path, a.String(), funName)
	}
}

func bindJson(out interface{}, c *gin.Context) error {
	err := c.ShouldBindJSON(out)
	if err != nil {
		log.Logger().Error("request err:", err)
		dieFail(err, c)
		return err
	}
	return nil
}

type HttpResponse struct {
	Code   int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Sign   int32  `protobuf:"varint,2,opt,name=sign,proto3" json:"sign,omitempty"`
	Msg    string `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	Result string `protobuf:"bytes,4,opt,name=msg,proto3" json:"result,omitempty"`
}

func dieFail(err error, c *gin.Context) {
	sentry.CaptureException(err)
	sentry.Flush(2 * time.Second)
	c.JSON(http.StatusOK, HttpResponse{
		Code:   400,
		Sign:   400,
		Msg:    err.Error(),
		Result: err.Error(),
	})
}

//健康检查响应函数
func healthResponse(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "health!",
	})
}
