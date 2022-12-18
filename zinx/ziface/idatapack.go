package ziface

/*
	封包、拆包，用于处理TCP黏包问题
	抽象层
*/

type IDataPack interface {
	// 获取包头的长度
	GetHeadLen() uint32
	// 封包方法
	Pack(msg IMessage) ([]byte, error)
	// 拆包方法
	UnPack([]byte) (IMessage, error)
}