package ziface

/*
	消息管理抽象层
*/
type IMsgHandle interface {
	DoMsgHandler(request IRequest)          //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑
}

//这里面有两个方法，`AddRouter()`就是添加一个msgId和一个路由关系到Apis中，那么`DoMsgHandler()`则是调用Router中具体`Handle()`等方法的接口。