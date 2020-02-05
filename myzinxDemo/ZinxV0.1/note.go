package main
/*
server  负责开1 获取一个TCP的Addr //2 监听服务器地址 //3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()

connect 负责 保存那个 就是conn, err := listenner.AcceptTCP() 那个通信的socket ,返回值是 *net.TCPConn
         还有存着一个 服回掉函数（ 拿到通信socket后，conn 要准备处理的业务） / 回掉函数还是也得拿着这个conn 才能干事情
                                                           type HandFunc func(*net.TCPConn, []byte, int) error
         connect没存请求写数据，请求数据靠存回掉函数

request  有点类似connect ，不同的是，读客户端发来的数据 放在connnect的stratReader方法里  ，  类似在需要把 客户端 请求的连接信息 和 请求的数据，放在一个叫Request的请求类里，这样的好处是我们可以从Request里得到全部客户端的请求信息，
			request结构体 理解为每次客户端的全部请求数据

		负责 保存 connect接口

router  服务端应用可以给Zinx框架配置当前链接的处理业务方法， 之前的Zinx-V0.2我们的Zinx框架处理链接请求的方法是固定的（就是那个connect的回掉函数），现在是可以自定义，并且有3种接口可以重写
		当然每个方法都有一个唯一的形参`IRequest`对象 ，也就是客户端请求过来的连接和请求数据，作为我们业务方法的输入数据。

		request 里有connect 和 []byte
		type IRouter interface{
			PreHandle(request IRequest)  //在处理conn业务之前的钩子方法
			Handle(request IRequest)	 //处理conn业务的方法
			PostHandle(request IRequest) //处理conn业务之后的钩子方法
		}

		type  server  struct     结构体里有  Router  ziface.IRouter   （//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务）
	    type Connection struct   结构体 也有   Router  ziface.IRouter  （//该连接的处理方法router）
		另外
		iserver 接口 多了个方法 AddRouter(router IRouter)  （//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用 ）

		个人： 相当于 将拿到的conn 和 用conn.read([]byte)读出来的已放在字节数组的数据，将这2个包装成request 。作为Router的输入数据，在connect的startReader方法，开启读客户端传来的数据那， 然后调用connect结构体自己的Router接口的方法 ，去处理这些
        以取代以前的调用connect 自己的回掉函数
		原来： server结构体和connect结构体里的  Router  ziface.IRouter 是有关联的 ，一开 传进server里，然后server 通过 NewConntion(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection传入connect结构体里，
		然后调用 c.Router.Handle(request)  --》 此里面request集合了conn的socket  和c.Conn.Read(buf)的buf ，

*/

////1 当然，外部类型也可以声明同名的字段或者方法，来覆盖内部类型的 ，不管我们如何同名覆盖，都不会影响内部类型，我们还可以通过访问内部类型来访问它的方法、属性字段等。
///2 对于内部类型的属性和方法访问上，我们可以用外部类型直接访问，也可以通过内部类型进行访问；但是我们为外部类型新增的方法属性字段，只能使用外部类型访问，因为内部类型没有这些。
//3 嵌入类型的强大，还体现在：如果内部类型实现了某个接口，那么外部类型也被认为实现了这个接口。我们稍微改造下例子看下。