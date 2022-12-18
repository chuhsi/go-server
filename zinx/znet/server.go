package znet

import (
	"fmt"
	"log"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名字
	Name string
	// 服务器版本
	IPVersion string
	// 服务器ip
	IP string
	// 服务器端口
	Port int
	// 多个路由
	// Router ziface.IRouter
	MsgsHandler ziface.IMsgs

	ConnManege ziface.IConnManager

	OnConnStart func(conn ziface.IConnetion)

	OnConnStop func(conn ziface.IConnetion)
}

// 初始化Server模块的方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgsHandler: NewMsgsHandle(),
		ConnManege:  NewConnManeger(),
	}
	return s
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnManege
}

// func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
// 	log.Println("[Conn Handle] CallbackToClient")
// 	if _, err := conn.Write(data[:cnt]); err != nil {
// 		log.Println("write back buf err", err)
// 		return errors.New("CallBackToClient err")
// 	}
// 	return nil
// }

// 启动服务器
func (s *Server) Start() {
	log.Printf("[Zinx] Server Name: %s, Listenner at IP: %s, Port:%d is starting ", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	log.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPackageSize: %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		// 开启消息队列工作池
		s.MsgsHandler.StartWorkerPool()
		// 1 获取一个TCP的地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Println("net.ResolveTCPAddr err: ", err)
			return
		}
		// 2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Println("net.ListenTCP err: ", err)
			return
		}
		fmt.Println("start Zinx server success", s.Name, "succeed Listenning...")

		var c_id uint32 = 0
		// 3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				log.Println("listenner.AcceptTCP err: ", err)
				continue
			}
			// 判断连接是否超过最大个数（4）
			if s.ConnManege.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("too many conns ...")
				fmt.Println("too many conns ...")
				fmt.Println("too many conns ...")
				conn.Close()
				continue
			}

			dealConn := NewConnetion(s, conn, c_id, s.MsgsHandler)
			c_id++
			// 启动当前业务处理
			go dealConn.Start()

			// 处理客户端连接业务（读写）
			/* go func() {
				for {
					buf := make([]byte, 512)
					// read
					cnt, err := conn.Read(buf)
					if err != nil {
						log.Println("conn.Read err: ", err)
						break
					}
					// writer
					if _, err := conn.Write(buf[:cnt]); err != nil {
						log.Println("conn.Write err: ", err)
						continue
					}
				}
			}()
			*/
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// 将一些服务器的资源，状态等一下信息进行回收
	fmt.Println("[STOP] Zinx server name", s.Name)
	s.ConnManege.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	s.Start()

	// todo 做服务器启动之后的额外业务

	// blocking status
	select {}
}

// 路由功能：给当前服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgsHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Succ! ")
}

func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnetion)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnetion)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnetion) {
	if s.OnConnStart != nil {
		fmt.Println("-----> Call OnConnStart() ...")
		s.OnConnStart(conn)
	}
}
func (s *Server) CallOnConnStop(conn ziface.IConnetion) {
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStop() ...")
		s.OnConnStop(conn)
	}
}
