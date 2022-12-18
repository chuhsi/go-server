package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

// 存储一切有关Zinx框架的全局参数，供其他模块使用
// 参数是可以通过zinx.json文件，由用户自己配置
type GlobalObj struct {
	// Server
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string
	// Zinx
	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
}

// 读取配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("/Users/max/go/src/zinx/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 定义一个全局对外的GlobalObj
var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "v0.4",
		TcpPort:        9999,
		Host:           "0.0.0.0",
		MaxConn:        100,
		MaxPackageSize: 4096,
		WorkerPoolSize: 5,
		MaxWorkerTaskLen: 10,
	}
	// 从conf/zinx.json
	GlobalObject.Reload()
}
