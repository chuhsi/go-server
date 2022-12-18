package ziface

import "net"

// 定义连接模块抽象层
type IConnetion interface {
	// 启动连接
	Start()
	// 停止连接
	Stop()
	// 获取当前的连接的绑定（socket）
	GetTCPConnetion() *net.TCPConn
	// 获取当前连接模块的连接ID
	GetConnID() uint32
	// 获取远程客户端TCP状态 IP port
	RemoteAddr() net.Addr
	// 发送数据给远程客户端
	SendMsg(uint32, []byte) error

	SetProperty(string, interface{})

	GetProperty(string) (interface{}, error)

	RemovaProperty(string)
}

// 定义一个处理连接业务的方法
// type HandleFunc func(*net.TCPConn, []byte, int) error
