package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
	// "log"
	// "fmt"
)

// 封包、拆包具体模块
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头的长度
func (dp *DataPack) GetHeadLen() uint32 {
	// Id uint32 + DateLen uint32 = 8字节
	return 8
}

// 封包方法
// func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
// 	// 创建一个存放bytes字节缓冲
// 	dataBuff := bytes.NewBuffer([]byte{})

// 	// 1 将DataLen写入dataBuff
// 	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMegDataLen()); err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	// 2 将MsgId写入dataBuff
// 	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMegId()); err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	// 3 将data写入dataBuff
// 	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMegData()); err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	return dataBuff.Bytes(), nil
// }

// // 拆包方法
// func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
// 	// 创建一个存放bytes字节缓冲
// 	dataBuff := bytes.NewReader(binaryData)

// 	msg := &Message{}

// 	// 1 从dataBuff读DataLen
// 	if err := binary.Read(dataBuff, binary.LittleEndian, msg.DateLen); err != nil {
// 		fmt.Println("从dataBuff读DataLen............")
// 		return nil, errors.New("从dataBuff读DataLen234567890")
// 	}
// 	// 2 从dataBuff读MsgId
// 	if err := binary.Read(dataBuff, binary.LittleEndian, msg.Id); err != nil {
// 		fmt.Println("从dataBuff读MsgId")
// 		return nil, err
// 	}
// 	// 判断datalen是否已经超出了我们允许的最大包的长度
// 	fmt.Println("判断datalen是否已经超出了我们允许的最大包的长度")
// 	if (utils.GlobalObject.MaxPackageSize > 0) && (msg.DateLen > utils.GlobalObject.MaxPackageSize) {
// 		return nil, errors.New("too large msg data recv!")
// 	}
// 	return msg, nil
// }

// Pack 封包方法(压缩数据)
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMegDataLen()); err != nil {
		return nil, err
	}

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMegId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMegData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//Unpack 拆包方法(解压数据)
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DateLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DateLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
