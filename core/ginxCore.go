package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mark-wby/ginx/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

//定义控制器的方法,返回任意类型
type  ControllerFunc func(ctx *gin.Context)interface{}


//将脚手架ginx中的控制器方法转化为gin的handlefunc
func ConvertHandlefunc(controllerFunc ControllerFunc) gin.HandlerFunc{
	return func(context *gin.Context) {
		//添加defer捕捉异常
		defer func() {
			if e :=recover();e!=nil{
				//捕捉到异常,断言是否是自定义的异常
				execption,ok := e.(GinxException)
				if ok {
					context.JSON(http.StatusOK,gin.H{
						"code":execption.Code,
						"status":false,
						"msg":execption.Message,
						"data":struct {}{},
					})
				}else {
					context.JSON(http.StatusOK,gin.H{
						"code":500,
						"status":false,
						"msg":"系统错误",
						"data":struct {}{},
					})
				}
			}
		}()

		//正常流程
		res := controllerFunc(context)
		context.JSON(http.StatusOK,gin.H{
			"code":200,
			"status":true,
			"msg":"调用成功",
			"data":res,
		})
	}
}

//核心实例
var GinCoreInstance *GinxCore;


//ginx脚手架的核心类
type GinxCore struct {
	*gin.Engine //gin引擎
	RouteGroup *gin.RouterGroup //路由分组
	GinxDb *gorm.DB //gorm
	RedisUtil *RedisUtil
	MqUtil *MqUtil
}


func NewGinxCore() *GinxCore {
	engine := gin.New()

	GinCoreInstance = &GinxCore{Engine: engine}
	return GinCoreInstance
}

//初始化数据库
func(this *GinxCore) InitDB(config *config.GinxDbConfig) *GinxCore{

	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=true",
		config.User, config.Password, config.Address, config.Database)

	fmt.Print("数据库链接字符串:"+dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err!=nil{
		log.Fatal("数据库链接失败")
	}

	//注入自定义插件,获取执行的sql
	db.Use(new(DbPlugin))

	this.GinxDb = db
	return this
}

//初始化mq
func(this *GinxCore) InitMq(mqConfig config.MqConfig) *GinxCore{
	this.MqUtil = NewMqUtil(mqConfig)
	return this
}

//初始化redis
func(this *GinxCore) InitRedis(redisConfig config.RedisConfig) *GinxCore{
	this.RedisUtil = NewRedisUtil(redisConfig)
	return this
}

//设置线上环境
func (this *GinxCore) SetGinReleaseMode() *GinxCore {
	gin.SetMode(gin.ReleaseMode)
	return this
}

//挂载全局中间件
func (this *GinxCore) Middleware(ginxMiddleware GinxMiddleware) *GinxCore{
	this.Engine.Use(ginxMiddleware.Handle)
	return this
}

//绑定路由(类似重载了gin的handle方法,只是为了在函数里面调用绑定方法)
func (this *GinxCore)Handle(httpMethod, relativePath string, handler ControllerFunc){
	//将控制器的绑定方法转为handlefunc
	handlefunc := ConvertHandlefunc(handler)
	//实际绑定路由的地方
	//判断是否使用Group
	if this.RouteGroup==nil{
		this.Engine.Handle(httpMethod,relativePath,handlefunc);
	}else {
		this.RouteGroup.Handle(httpMethod,relativePath,handlefunc);
	}
}

//路由分组
func (this *GinxCore) Group(group string) *GinxCore {
	//绑定路由分组
	this.RouteGroup =this.Engine.Group(group)
	return this
}

//绑定控制器
func (this *GinxCore) Build(controllers ...GinxController) *GinxCore{
	for _,controller := range controllers{
		controller.Bind(this)
	}
	return this
}

//启动函数
func (this *GinxCore) Start(port string){
	//注册404路由不存在处理
	this.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusOK,gin.H{
			"status":false,
			"msg":"请求路径不存在",
			"data":struct {}{},
			"code":404,
		})
	})

	//注册请求方法不存在处理
	this.NoMethod(func(context *gin.Context) {
		context.JSON(http.StatusOK,gin.H{
			"status":false,
			"msg":"请求方式不存在",
			"data": struct {}{},
			"code":405,
		})
	})
	err := this.Engine.Run(port)
	if err != nil {
		log.Fatal("启动失败:",err.Error())
	}
}

