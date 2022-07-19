package custom

type GinxException struct {
	Code int //错误码
	Message string //错误消息
}

func NewGinxException(code int, message string) *GinxException {
	return &GinxException{Code: code, Message: message}
}
