package znet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

/*
	连接模块
*/
type Connection struct {
	TCPServer ziface.IServer

	// 当前连接的Socket TCP
	Conn *net.TCPConn
	// 当前连接ID
	ConnID uint32
	// 当前连接状态
	IsClosed bool
	// 当前连接绑定的业务处理方法API
	// HandlerAPI ziface.HandleFunc
	// 当前退出的连接
	ExitChan chan bool

	//
	// Router ziface.IRouter
	MsgsHandler ziface.IMsgs

	// 无缓冲管道，用于读写通信
	MsgChan chan []byte

	property map[string]interface{}

	propertyLock sync.Mutex
}

// 初始化连接模块的方法
func NewConnetion(server ziface.IServer, conn *net.TCPConn, connID uint32, router ziface.IMsgs) *Connection {
	c := &Connection{
		TCPServer:   server,
		Conn:        conn,
		ConnID:      connID,
		MsgsHandler: router,
		IsClosed:    false,
		ExitChan:    make(chan bool, 1),
		MsgChan:     make(chan []byte),
		property:    make(map[string]interface{}),
	}
	// 添加到连接管理中
	c.TCPServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())
	for {
		select {
		case data := <-c.MsgChan:
			// 把数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("Send data err", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时也要退出
			return
		}
	}
}

func (c *Connection) StartReader() {
	log.Println("Reader Goroutine is running ... ")
	defer fmt.Println("[Reader is exit]", "ConnID = ", c.ConnID, "[remote addr is ]", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	log.Println("c.Conn.Read err ", err)
		// 	continue
		// }

		// 创建封包、拆包的对象
		dp := NewDataPack()
		// 读取客户端的Msg head 二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnetion(), headData); err != nil {
			log.Fatalln("read msg head err ", err)
			break
		}
		// 拆包 得到包头信息，ID，Len
		msg, err := dp.Unpack(headData)
		if err != nil {
			log.Fatalln("unpack err ", err)
			break
		}
		// 根据Len，读取Data
		var data []byte
		if msg.GetMegDataLen() > 0 {
			data = make([]byte, msg.GetMegDataLen())
			if _, err := io.ReadFull(c.GetTCPConnetion(), data); err != nil {
				log.Fatalln("read msg head err ", err)
				break
			}
		}
		msg.SetMegData(data)
		req := &Request{
			Conn: c,
			Msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			go c.MsgsHandler.SendMsgToTaskQue(req)
		} else {
			go c.MsgsHandler.DoMsgHandle(req)
		}

		// go c.MsgsHandler.DoMsgHandle(req)
		// go func(request ziface.IRequest) {
		// 	c.MsgsHandler.DoMsgHandle(req)
		// }(req)
		/*
			调用当前连接所绑定的HandlerAPI
			if err := c.HandlerAPI(c.Conn, buf, cnt); err != nil {
				log.Println("ConnID ",c.ConnID, "HandlerAPI is err", err)
				break
			}
		*/
	}
}

func (c *Connection) Start() {
	log.Println("Conn start ... ConnID = ", c.ConnID)
	// 读数据业务
	go c.StartReader()
	// 写数据业务
	go c.StartWriter()
	// 连接之前所做的工作
	c.TCPServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	log.Println("Conn stop ... ConnID = ", c.ConnID)
	if c.IsClosed {
		return
	}
	c.IsClosed = true

	c.TCPServer.CallOnConnStop(c)

	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	c.TCPServer.GetConnMgr().Remove(c)

	// 关闭管道
	close(c.ExitChan)
	close(c.MsgChan)
}

func (c *Connection) GetTCPConnetion() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msdId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("Connection closed when send msg")
	}
	// 将数据进行封包 Len|Id|Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msdId, data))
	if err != nil {
		log.Println("dp.Pack(NewMessage(msdId, data)) ", err)
		return errors.New("dp.Pack(NewMessage(msdId, data))")
	}
	// 将数据发送给MsgChan
	c.MsgChan <-binaryMsg
	// if _, err := c.Conn.Write(binaryMsg); err != nil {
	// 	log.Println("c.Conn.Write(binaryMsg) ", err)
	// 	return errors.New("c.Conn.Write(binaryMsg) ")
	// }
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if val, ok := c.property[key]; ok {
		return val, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemovaProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
