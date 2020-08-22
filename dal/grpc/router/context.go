package router

import (
	pb "dj-api/proto"
	"math"
)

type HandlerFunc func(*Context) (*pb.EMRsp, error)

const abortIndex int8 = math.MaxInt8 / 2

// HandlersChain defines a HandlerFunc array.
type HandlersChain []HandlerFunc

//定义的上下文
type Context struct {
	Req      *pb.EMReq
	Rsp      *pb.EMRsp
	handlers HandlersChain
	index    int8
}

//模拟的调用堆栈
func (c *Context) Next() (rsp *pb.EMRsp, err error) {
	c.index++
	for c.index < int8(len(c.handlers)) {
		//按顺序执行HandlersChain内的函数
		//如果函数内无c.Next()方法调用则函数顺序执行完
		//如果函数内有c.Next()方法调用则代码执行到c.Next()方法处压栈，等待后面的函数执行完在回来执行c.Next()后的命令
		rsp, err = c.handlers[c.index](c)
		if err != nil {
			c.Abort()
		}
		c.index++
	}
	return
}

func (c *Context) Abort() {
	c.index = abortIndex
}
func (c *Context) reset() {
	c.handlers = nil
	c.index = -1
}
