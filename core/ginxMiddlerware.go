package core

import "github.com/gin-gonic/gin"

//框架的中间件接口
type GinxMiddleware interface {

	Handle(context *gin.Context);//处理流程
}
