package config

//数据库配置模型
type GinxDbConfig struct {
	Driver   string
	Address  string
	Database string
	User     string
	Password string
}

//mq配置
type MqConfig struct {
	Host string
	Port int
	UserName string
	Password string
}

//redis配置
type RedisConfig struct {
	Host string
	Port string
	Password string
	Db int
}
