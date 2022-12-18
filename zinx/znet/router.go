package znet

import "zinx/ziface"

// 实现Router时，先嵌入BaseRouter基类，然后根据具体需求进行方法重写
type BaseRouter struct {}

//处理Conn业务之前的方法
func (b *BaseRouter) PreHandle(request ziface.IRequest){

}
//处理Conn业务的主方法
func (b *BaseRouter) Handle(request ziface.IRequest){

}
//处理Conn业务之的方法
func (b *BaseRouter) PostHandle(request ziface.IRequest){

}