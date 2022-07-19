package middleware

import (
	"bytes"
	"encoding/json"
	"ginx/core"
	"ginx/custom"
	"ginx/util"
	"github.com/gin-gonic/gin"
)

//日志中间件
type RequestLogMiddleware struct {

}


//中间件处理函数
func(this *RequestLogMiddleware) Handle(context *gin.Context){
	//替换自定义的response(可以存储响应内容)
	blw := &custom.CustomResponseWrite{
		Body:           bytes.NewBufferString(""),
		ResponseWriter: context.Writer,
		LogUtil:        util.NewLoggerUtil(),
	}

	//解析任何请求方式的请求参数,塞入结构体中
	//解析get请求和post请求的form参数
	//request.ParseForm()
	context.Request.ParseMultipartForm(128)
	// 获取请求实体长度
	contentLength := context.Request.ContentLength
	body := make([]byte, contentLength)
	//获取请求体数据
	context.Request.Body.Read(body)
	//定义map结构接受数据
	event := make(map[string]interface{},0)
	//将json数据解析成map
	json.Unmarshal(body, &event)
	//由于其他不是json请求获取到参数不能满足map,需要进行转化
	tmpData := make(map[string]interface{},0)
	for k,v := range context.Request.Form{
		if len(v)>1{
			tmpData[k] = v
		}else {
			tmpData[k] = v[0]
		}

	}
	blw.RequestParam = util.MergeMap(tmpData,event)

	//关键地方在这里,将上线文中的writer对象替换成自定义的writer对象
	context.Writer = blw

	//j
	core.RequestContext = context

	context.Set("custom",blw)

	context.Next()


	//记录日志提交到mq
	//将请求返回值变成map
	responseData := make(map[string]interface{},10)

	//将返回值解析成map
	json.Unmarshal(blw.Body.Bytes(),&responseData)

	//fmt.Println(string(write.Body.Bytes()))

	//msg :=map[string]interface{}{
	//	"requestPath":context.Request.URL.Path,
	//	"requestParam":blw.RequestParam,
	//	"requestResponse":responseData,
	//	"requestSqlLog":blw.LogUtil.GetSqlLog(),
	//	"requestLog":blw.LogUtil.GetCustomLog(),
	//	"requestProjectName":"ceshi",
	//	"createdAt":context.Request.Header.Get("startTime"),
	//	"updatedAt":time.Now().UnixNano(),
	//	"requestId":fmt.Sprintf("%x",md5.Sum([]byte(context.Request.Header.Get("plusTraceId")))),
	//	"uniqueTraceId":context.Request.Header.Get("plusTraceId"),
	//}
	//res,_ := json.Marshal(msg)
	//
	//core.GinCoreInstance.MqUtil.PushMsg(string(res),"logExchange","test")
}