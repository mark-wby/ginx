package core


//框架的控制器接口
type GinxController interface {
	Bind(ginxCore *GinxCore)//绑定路由

}
