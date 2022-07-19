package main

import (
	"encoding/json"
	"github.com/mark-wby/ginx/config"
	"github.com/mark-wby/ginx/core"
	"github.com/mark-wby/ginx/demo/controller"
	"log"
	"os"
)

//测试案例
func main(){
	//解析数据库参数
	config := jsonDecodeDbConfig();
	core.NewGinxCore().
		InitDB(config).
		Build(&controller.IndexController{}).
		Start(":8888")

}

//解析数据库配置
func jsonDecodeDbConfig()*config.GinxDbConfig {
	file, err := os.Open("config/dbConfig.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config := &config.GinxDbConfig{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
	return config
}