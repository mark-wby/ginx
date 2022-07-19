package util

import (
	"bytes"
	"encoding/json"
	"ginx/custom"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)


//将map转换为指定的struct
func MapToStruct(m map[string]interface{},str interface{})interface{}{
	strType := reflect.ValueOf(str)
	if strType.Kind() == reflect.Ptr {
		strType = strType.Elem()
	}
	if strType.Kind() != reflect.Struct {
		panic("str must be a struct type")
	}

	//遍历字典
	for k,v := range m {
		for i := 0; i < strType.NumField(); i++ {
			if strType.Type().Field(i).Name == k {
				//结构体的字段和字典的键匹配,判断结构体的字段是否允许赋值
				if strType.Field(i).CanSet() {
					strType.Field(i).Set(reflect.ValueOf(v))
					break
				}
			}
		}
	}

	return str
}

//将时间戳转为时间格式(时分秒格式为:2006-01-02 15:04:05)
func TimestampToDateFormat(timestamp int64,format string) string{
	tm := time.Unix(timestamp,0)
	return tm.Format(format)
}

//将时间格式专为时间戳
func DateFormatTOTimestamp(format string) int64{
	tt,_ := time.ParseInLocation("2006-01-02 15:04:05",format,time.Local)
	return tt.Unix()
}

//获取当前时间戳
func GetNowTimestamp() int64 {
	return time.Now().Local().Unix()
}

//合并map
func MergeMap(mapDatas ...map[string]interface{}) map[string]interface{}{
	tmpMap := make(map[string]interface{},0)

	for _,mapData := range mapDatas{
		for k,v := range mapData{
			tmpMap[k] = v
		}
	}

	return tmpMap
}

//封装httpget请求
func HttpGetRequest(url string,headers map[string]string,res interface{}) interface{}{
	client := &http.Client{}
	//构建自定义http请求
	req, _ := http.NewRequest("GET", url, nil)
	//塞入请求头
	req.Header.Add("Content-Type","application/json")
	for k,v := range headers{
		req.Header.Add(k, v)
	}
	//执行请求
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err!=nil{
		panic(custom.NewGinxException(500,err.Error()))
	}
	//读取请求数据
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		panic(custom.NewGinxException(500,err.Error()))
	}

	err = json.Unmarshal(body,res)
	if err != nil {
		panic(custom.NewGinxException(500,"解析失败"))
	}
	return res
}

//封装httppost请求
func HttpPostRequest(url string,headers map[string]string,params map[string]interface{},res interface{}) interface{}{
	client := &http.Client{}
	data,_ := json.Marshal(params)
	bodyReader := bytes.NewReader(data)
	//构建自定义http请求
	req, _ := http.NewRequest("POST", url, bodyReader)
	//塞入请求头
	req.Header.Add("Content-Type","application/json")
	for k,v := range headers{
		req.Header.Add(k, v)
	}
	//执行请求
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err!=nil{
		panic(custom.NewGinxException(500,err.Error()))
	}
	//读取请求数据
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		panic(custom.NewGinxException(500,err.Error()))
	}

	err = json.Unmarshal(body,res)
	if err != nil {
		panic(custom.NewGinxException(500,"解析失败"))
	}
	return res
}
