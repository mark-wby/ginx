package custom

import (
	"bytes"
	"github.com/mark-wby/ginx/util"
	"github.com/gin-gonic/gin"
)

type CustomResponseWrite struct {
	gin.ResponseWriter
	Body *bytes.Buffer
	LogUtil *util.LoggerUtil
	RequestParam map[string]interface{}
}

func (w CustomResponseWrite) Write(b []byte) (int,error){
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWrite) WriteString(s string) (int,error){
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
