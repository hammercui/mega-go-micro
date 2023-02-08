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
package gin

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/hammercui/mega-go-micro/v2/log"
	"github.com/micro/go-micro/v2/server"
	"net/http"
	"reflect"
	"strings"
	"time"
)

//Server is a simple micro server abstraction
//GinServer 实现Server接口
type GinServer struct {
	ginRouter *gin.Engine
	basePath  string
	options server.Options
	megaGinOptions *megaGinServerOptions
}

func (p *GinServer) SetBasePath(basePath string) {
	p.basePath = basePath
}

func (p *GinServer) NewSubscriber(s string, i interface{}, option ...server.SubscriberOption) server.Subscriber {
	log.Logger().Infof("NewSubscriber func")
	return nil
}

func (p *GinServer) Subscribe(subscriber server.Subscriber) error {
	log.Logger().Infof("Subscribe func")
	return nil
}

func (p *GinServer) Start() error {
	log.Logger().Infof("Start func")
	return nil
}

func (p *GinServer) Stop() error {
	log.Logger().Infof("Stop func")
	return nil
}

func (p *GinServer) String() string {
	log.Logger().Infof("String func")
	return "custom GinServer"
}
func (p *GinServer) Init(option ...server.Option) error {
	log.Logger().Infof("GinServer Init")
	return nil
}

func (p *GinServer) Options() server.Options {
	return p.options
}

func (p *GinServer) Server() *GinServer {
	return p
}

func (p *GinServer) Gin() *gin.Engine {
	return p.ginRouter
}

func (p *GinServer) Handle(handler server.Handler) error {
	return nil
}

func (p *GinServer) NewHandler(i interface{}, opts ...server.HandlerOption) server.Handler {
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
		p.registerEndPointsV2Imp(funName, method, path, i)

	}
	return nil
}

//注册endpoint实现
func (p *GinServer) registerEndPointsV2Imp(funName string, method string, path string, service interface{}) {
	fullPath := fmt.Sprintf("%s%s", p.basePath, path)
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
			if erro := p.bindJson(req.Interface(), c); erro != nil {
				return
			}
			parameters[2] = reflect.ValueOf(req.Interface())

			//出参
			respPtr := m.Type.In(3)
			var resp = reflect.New(respPtr.Elem())
			parameters[3] = reflect.ValueOf(resp.Interface())

			//执行函数
			errorResult := m.Func.Call(parameters)
			p.handleEndpointResult(c,errorResult,resp)
		}
		switch method {
		case "POST":
			p.ginRouter.POST(fullPath, a)
		case "GET":
			p.ginRouter.GET(fullPath, a)
		}

	} else {
		log.Logger().Errorf("注册service路由失败(post):%s,类%s ,未实现方法%s", fullPath, a.String(), funName)
	}
}

//处理函数执行结果
func (p *GinServer) handleEndpointResult(c *gin.Context,errorResult []reflect.Value,resp reflect.Value)  {
	//返回结果
	for _, err := range errorResult {
		if err.Interface() == nil {
			var body = make(gin.H)
			for _, f := range p.megaGinOptions.responseFields {
				switch f.FieldType {
				case "string":
					body[f.Name] = "ok"
				case "int":
					body[f.Name] = p.megaGinOptions.responseSuccessCode
				case "bool":
					body[f.Name] = true
				case "interface":
					body[f.Name] = resp.Interface()
				}
			}
			c.JSON(http.StatusOK, body)
		} else {
			p.dieFail(err.Interface().(error), c)
		}
	}
}

func (p *GinServer) bindJson(out interface{}, c *gin.Context) error {
	err := c.ShouldBindJSON(out)
	if err != nil {
		log.Logger().Error("request err:", err)
		p.dieFail(err, c)
		return err
	}
	return nil
}

func (p *GinServer) dieFail(err error, c *gin.Context) {
	sentry.CaptureException(err)
	sentry.Flush(2 * time.Second)
	message := err.Error()
	var body = make(gin.H)
	for _, item := range p.megaGinOptions.responseFields {
		switch item.FieldType {
		case "string":
			body[item.Name] = message
		case "int":
			body[item.Name] = p.megaGinOptions.responseFailCode
		case "bool":
			body[item.Name] = false
		}
	}
	c.JSON(http.StatusBadRequest, body)
}

//健康检查响应函数
func (p *GinServer) healthResponse(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "health!",
	})
}

func (p *GinServer) pongResponse(context *gin.Context) {
	var body = make(gin.H)
	for _, item := range p.megaGinOptions.responseFields {
		switch item.FieldType {
		case "string":
			body[item.Name] = "pong!"
		case "bool":
			body[item.Name] = true
		}
	}
	context.JSON(http.StatusOK, body)
}

func (p *GinServer) Apply(options ...MegaGinServerOption) *GinServer {
	for _, opt := range options{
		opt.apply(p.megaGinOptions)
	}
	return p
}
//实例化megaGinServer
func NewMegaGinServer(middlewares ...gin.HandlerFunc) *GinServer {
	gin.DisableConsoleColor()
	//gin设置模式
	gin.SetMode(gin.DebugMode)
	//初始化路由
	ginRouter := gin.New()
	//注册通用中间件
	for _, item := range middlewares {
		ginRouter.Use(item)
	}
	ginServer := &GinServer{
		ginRouter: ginRouter,
		megaGinOptions: defaultMegaGinServerOptions(),
	}
	//健康检查
	ginRouter.GET("", ginServer.healthResponse)
	ginRouter.GET("/actuator", ginServer.healthResponse)
	ginRouter.GET("/actuator/health", ginServer.healthResponse)
	//心跳
	ginRouter.GET("/ping", ginServer.pongResponse)

	return ginServer
}