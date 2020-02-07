package znet

import (
	"fmt"
	"pororo.com/zinx/utils"
	"pororo.com/zinx/ziface"
	"strconv"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法的map属性
	//一个woker池（多个worker）， 一个消息队列，6个worker 6个消息队列（6个切片）
	WorkerPoolSize uint32                 //业务工作Worker池的数量
	TaskQueue      []chan ziface.IRequest //Worker负责取任务的消息队列 注意是channel切片
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
		//`WokerPoolSize`:作为工作池的数量，因为TaskQueue中的每个队列应该是和一个Worker对应的，所以我们在创建TaskQueue中队列数量要和Worker的数量一致。
		//`TaskQueue`是一个Request请求信息的channel集合。用来缓冲提供worker调用的Request请求信息，worker会从对应的队列中获取客户端的请求数据并且处理掉。
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

// 增加工作池 消息队列功能
//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen) // 代表chan容量 是 1024（最大值）
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
	// chan1(1024容量)  chan2 chan3 chan4 chan5 chan6 chan7 chan8 chan9
}

/*`SendMsgToTaskQueue()`作为工作池的数据入口，这里面采用的是轮询的分配机制，因为不同链接信息都会调用这个入口，那么到底应该由哪个worker处理该链接的请求处理，
 整理用的是一个简单的求模运算。用余数和workerID的匹配来进行分配。
​ 最终将request请求数据发送给对应worker的TaskQueue，那么对应的worker的Goroutine就会处理该链接请求了。*/

//将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}
