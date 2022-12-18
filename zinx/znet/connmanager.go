package znet

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"zinx/ziface"
)

type ConnManeger struct {
	Conns    map[uint32]ziface.IConnetion
	ConnLock sync.RWMutex
}

func NewConnManeger() *ConnManeger {
	return &ConnManeger{
		Conns: make(map[uint32]ziface.IConnetion),
	}
}

func (cm *ConnManeger) Add(conn ziface.IConnetion) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()
	// 将conn加入到ConnManeger中
	cm.Conns[conn.GetConnID()] = conn
	log.Println("connection add to ConnManeger successfully: conn num = ", cm.Len())
}

func (cm *ConnManeger) Remove(conn ziface.IConnetion) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	delete(cm.Conns, conn.GetConnID())
	log.Println("connID = ", conn.GetConnID(), "remove from ConnManeger successfully: conn num = ", cm.Len())
}

func (cm *ConnManeger) Get(connID uint32) (ziface.IConnetion, error) {
	cm.ConnLock.RLock()
	defer cm.ConnLock.RUnlock()

	if conn, ok := cm.Conns[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!!!")
	}
}

func (cm *ConnManeger) Len() int {
	return len(cm.Conns)
}

func (cm *ConnManeger) ClearConn() {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	for connID, conn := range cm.Conns {
		conn.Stop()
		delete(cm.Conns, connID)
	}
	fmt.Printf("Clear all connections succ!! conn num = %d", cm.Len())
}
