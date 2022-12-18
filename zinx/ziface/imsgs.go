package ziface

// 消息管理抽象层

type IMsgs interface {
	// 调度对应的路由器
	DoMsgHandle(request IRequest)
	// 添加路由器
	AddRouter (uint32, IRouter)
	
	StartWorkerPool()

	StartOneWorker(uint32, chan IRequest)
	
	SendMsgToTaskQue(IRequest)
}