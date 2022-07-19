package util

import (
	"fmt"
	"github.com/mark-wby/ginx/config"
	"github.com/streadway/amqp"
	"log"
)



type MqUtil struct {
	MqConn *amqp.Connection
}

func NewMqUtil(config config.MqConfig) *MqUtil {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",config.Host,config.UserName,
		config.Password,config.Port)
	fmt.Println(dsn)
	conn,err := amqp.Dial(dsn)
	if err !=nil{
		fmt.Println("mq链接失败")
		log.Fatal(err)
	}
	return &MqUtil{MqConn: conn}
}


//推送消息到mq
func (this *MqUtil) PushMsg(body string,exchange string,logKey string) error{
	channel,err := this.MqConn.Channel()
	if err!= nil{
		log.Fatal(err)
	}
	defer channel.Close()

	err = channel.Publish(exchange,logKey,false,false,
		amqp.Publishing{
			ContentType:"text/plain",
			Body:[]byte(body),
		})
	return err
}

