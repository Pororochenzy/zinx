package utils

import (
	"encoding/json"
	"io/ioutil"
	"pororo.com/zinx/ziface"
)

//  随着架构逐步的变大，参数就会越来越多，为了省去我们后续大频率修改参数的麻烦，接下来Zinx需要做一个加载配置的模块，和一个全局获取Zinx参数的对象。
/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数也可以通过 用户根据 zinx.json来配置
*/
//全局定义了一个`GlobalObject`对象，目的就是让其他模块都能访问到里面的参数。
/*
	定义一个全局的对象
*/
var GlobalObject *GlobalObj

type GlobalObj struct {
	TcpServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称
	Version   string         //当前Zinx版本号

	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
}

/*
	提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//从配置文件中加载一些用户配置的参数
	// GlobalObject.Reload()
}

//读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json") // 注意： 不是看，哪个函数调用 的 ，就看他的那个函数所在文件的位置， 而是哪个目录 运行go run Server 命令
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
