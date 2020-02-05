package znet

import "pororo.com/zinx/ziface"

type Request struct {
	conn ziface.IConnection //已经和客户端建立好的 链接
	//，现在是用一个`[]byte`来接受全部数据，又没有长度，又没有消息类型，这不科学。怎么办呢？我们现在就要自定义一种消息类型，把全部的消息都放在这种消息类型里。
	// data []byte //客户端请求的数据
	msg ziface.IMessage 	//客户端请求的数据
}
//获取请求连接信息
func(r *Request) GetConnection() ziface.IConnection {
	return r.conn
}
//获取请求消息的数据
func(r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}