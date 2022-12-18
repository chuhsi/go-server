package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"zinx/znet"
)

/*
	模拟客户端
*/
func main() {
	fmt.Println("Client start ...")
	time.Sleep(1 * time.Second)

	// 1 连接客户端，得到一个连接Conn
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Println("net.Dial err: ", err)
		return
	}
	fmt.Println("conn start ...")

	// 2 写入数据
	for {
		// _, err := conn.Write([]byte("Hello Zinx"))
		// if err != nil {
		// 	log.Println("conn.Write err: ", err)
		// 	return
		// }
		// buf := make([]byte, 512)
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	log.Println("conn.Read err: ", err)
		// 	return
		// }
		// fmt.Printf("Server call back: %s count = %d\n", buf, n)

		// 发送封包的Message消息
		dp := znet.NewDataPack()
		binary, err := dp.Pack(znet.NewMessage(1, []byte("zinx v0.5 client test msg")))
		fmt.Println("Pack start ...")
		if err != nil {
			fmt.Println("dp.Pack err", err)
			return
		}
		if _, err := conn.Write(binary); err != nil {
			fmt.Println("conn.Write err", err)
			return
		}
		fmt.Println("Write start ...")

		//1 服务器返回的数据进行拆包处理
		// binaryHead := make([]byte, dp.GetHeadLen())
		// if _, err := io.ReadFull(conn, binaryHead); err != nil {
		// 	fmt.Println("io.ReadFull(conn,binaryHead)", err)
		// 	break
		// }
		// //2 先读取流中 head部分， 得到ID，Len
		// msgHead, err := dp.UnPack(binaryHead)
		// if err != nil {
		// 	fmt.Println("dp.UnPack(binaryHead) err", err)
		// 	return
		// }
		headData := make([]byte, dp.GetHeadLen())
		fmt.Println("make([]byte, dp.GetHeadLen()) start ...")
		if _, err = io.ReadFull(conn, headData); err != nil {
			fmt.Println(headData)
			log.Fatalln("read msg head err ", err)
			break
		}
		fmt.Println("io.ReadFull(conn, headData) start ...")
		// 拆包 得到包头信息，ID，Len
		data, err := dp.Unpack(headData) 
		if err != nil {
			log.Fatalln("unpack err ", err)
			break
		}
		fmt.Println("dp.UnPack(headData) start ...")
		if data.GetMegDataLen() > 0 {
			msg := data.(*znet.Message)
			msg.Data = make([]byte, msg.GetMegDataLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("io.ReadFull(conn,msg.Data) err", err)
				return
			}
			fmt.Println("------>Recv Server Msg: ID=", msg.Id, "Len=", msg.DateLen, "data=", string(msg.Data))
		}
		//3 再根据Len，读取data部分
		time.Sleep(1 * time.Second)
	}
}
