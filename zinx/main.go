package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

func DoConnBegin(conn ziface.IConnetion) {
	fmt.Println("DoConnBegin is called")
	err := conn.SendMsg(202, []byte("DoConnBegin"))
	if err != nil {
		fmt.Println(err)
	}
	conn.SetProperty("name", "max")
}

func DoConnAfterClosed(conn ziface.IConnetion) {
	fmt.Println("DoConnAfterClosed is called")
	fmt.Println("conn ID = ", conn.GetConnID(), "is done")
	fmt.Println(conn.GetProperty("name"))
}

func main() {
	// 1 创建一个Server句柄，使用Zinx的api
	// s := znet.NewServer("zinx 0.1")
	// s := znet.NewServer("zinx 0.2")
	// s := znet.NewServer("zinx 0.3")
	// s := znet.NewServer()//0.5
	// s := znet.NewServer()//0.6
	// s := znet.NewServer() //0.7
	s := znet.NewServer() //0.8
	// 2 注册hook回调函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnAfterClosed)
	// 3 添加Router
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	// 4 启动Server
	s.Serve()
}

type PingRouter struct {
	znet.BaseRouter
}

func (*PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	// request.GetConnetion().GetTCPConnetion().Write([]byte("handle ping ...\n"))
	request.GetConnetion().SendMsg(200, []byte("ping...ping...ping"))
}

type HelloRouter struct {
	znet.BaseRouter
}

func (*HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	// request.GetConnetion().GetTCPConnetion().Write([]byte("handle ping ...\n"))
	request.GetConnetion().SendMsg(200, []byte("hello...hello...hello"))
}

// test pre
// func (*PingRouter) PreHandle(request ziface.IRequest) {
// 	fmt.Println("Call Router PreHandle")
// 	request.GetConnetion().GetTCPConnetion().Write([]byte("before ping ...\n"))
// }
// test handle
// test post
// func (*PingRouter) PostHandle(request ziface.IRequest) {
// 	fmt.Println("Call Router PostHandle")
// 	request.GetConnetion().GetTCPConnetion().Write([]byte("after ping ...\n"))
// }
