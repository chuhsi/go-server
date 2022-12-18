package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

// 消息处理模块的实现

type MsgsHandle struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 消息队列
	TaskQue []chan ziface.IRequest
	// worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgsHandle() *MsgsHandle {
	return &MsgsHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQue:        make([]chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// 启动一个Worker工作池(只能启动一个)
func (m *MsgsHandle) StartWorkerPool() {
	// 根据WorkerPoolSize 分别开启Worker 每个worker用一个go来开启
	var i uint32
	for i = 0; i < m.WorkerPoolSize; i++ {
		m.TaskQue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前Worker， 阻塞消息从channel传递过来
		go m.StartOneWorker(i, m.TaskQue[i])
	}
}

// 启动一个Worker工作流程
func (m *MsgsHandle) StartOneWorker(workerID uint32, taskQue chan ziface.IRequest) {
	fmt.Println("WorkerId = ", workerID, "is started")
	for {
		select {
		case request := <-taskQue:
			m.DoMsgHandle(request)
		}
	}
}

func (m *MsgsHandle) DoMsgHandle(request ziface.IRequest) {
	handle, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), "is not FUND! need register")
	}
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

func (m *MsgsHandle) SendMsgToTaskQue(request ziface.IRequest) {
	// 将消息平均分配给Worker
	workerID := request.GetConnetion().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnetion().GetConnID(), "request MsgID = ", request.GetMsgId(), "to WorkerID = ", workerID)
	// 将消息发送给对应的Worker的TaskQue
	m.TaskQue[workerID] <-request
}

func (m *MsgsHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1 判断 当前msg绑定的API处理方法是否存储
	if _, ok := m.Apis[msgID]; ok {
		panic("repeqt api, msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2 添加msgID和API的关系
	m.Apis[msgID] = router
	fmt.Println("add API MsgID = ", msgID, "success")
}
