package znet

import "zinx/ziface"

type Request struct {
	// 已经和客户端建立好的连接
	Conn ziface.IConnetion
	// 客户端请求的数据
	// Data []byte
	Msg ziface.IMessage
}

func (r *Request) GetConnetion() ziface.IConnetion {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Msg.GetMegData()
}

func (r *Request) GetMsgId() uint32 {
	return r.Msg.GetMegId()
}
