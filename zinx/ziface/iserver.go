package ziface

// 定义服务接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()

	// 路由功能：给当前服务注册一个路由方法，供客户端的连接处理使用
	AddRouter(uint32, IRouter)

	GetConnMgr() IConnManager

	SetOnConnStart(func(connection IConnetion))

	SetOnConnStop(func(connection IConnetion))

	CallOnConnStart(connection IConnetion)
	CallOnConnStop(connection IConnetion)
}
