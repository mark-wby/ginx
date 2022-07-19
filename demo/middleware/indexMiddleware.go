package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type IndexMiddleware struct {

}

func (this *IndexMiddleware) Handle(ctx *gin.Context){
	fmt.Println("index中间件")
	ctx.Next()
}
