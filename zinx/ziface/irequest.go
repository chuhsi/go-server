package ziface

/*
	把客户端请求的连接信息和请求的数据包装到一个Request中
*/

type IRequest interface {
	// 获得当前连接
	GetConnetion() IConnetion
	// 获得当前连接的数据
	GetData() []byte

	GetMsgId() uint32
}