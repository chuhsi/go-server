package ziface

type IConnManager interface {

	Add(conn IConnetion)

	Remove(conn IConnetion)

	Get(connID uint32) (IConnetion, error)

	Len() int
	
	ClearConn()
}