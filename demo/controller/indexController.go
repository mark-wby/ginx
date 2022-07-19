package controller

import (
	"github.com/mark-wby/ginx/core"
	"github.com/gin-gonic/gin"
)

type IndexController struct {

}

func NewIndexController() *IndexController {
	return &IndexController{}
}

type vtool struct {
	Code float64
	Msg string
	Msec float64
	Time float64
	Data interface{}
}
//测试方法
func (this *IndexController) ceshi(context *gin.Context) interface{}{

	return "ceshi"
}

func (this *IndexController) Bind(ginxCore *core.GinxCore){
	//绑定路由
	ginxCore.Handle("GET","/",this.ceshi);

}
