package ziface

// 将请求消息封装到Message中，定义Message抽象接口
type IMessage interface {
	// getter
	GetMegId() uint32
	GetMegDataLen() uint32
	GetMegData() []byte
	// setter
	SetMegId(uint32)
	SetMegDataLen(uint32)
	SetMegData([]byte)
}