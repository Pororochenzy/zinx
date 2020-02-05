package ziface

import "net"

//定义连接接口 （通信的那个 sockert， 也就是 conn, err := listenner.AcceptTCP()）
type IConnection interface {
	//启动连接，让当前连接开始工作
	Start()
	//停止连接，结束当前连接状态M
	Stop()
	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	//获取当前连接ID
	GetConnID() uint32
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
	//直接将Message数据发送数据给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error
}

//定义一个统一处理链接业务的接口 （一个通信的sokcert 和 客户端请求的数据字节，请求数据长度。 作用让这条socket处理业务，例如做重复读的操作）
type HandFunc func(*net.TCPConn, []byte, int) error
