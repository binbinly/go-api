package router

import (
	pb "dj-api/proto"
	"sync"
)

type Engine struct {
	tree map[string]HandlersChain // tree为了简化做成了map路由路径完全匹配
	Group
	pool sync.Pool // 正常情况存在大量的上下文切换，所以使用一个临时对象存储
}

func NewEngine() *Engine {
	engine := &Engine{
		Group: Group{
			Handlers: nil,
		},
		tree: make(map[string]HandlersChain),
	}
	engine.Group.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

func (engine *Engine) allocateContext() *Context {
	return &Context{}
}

//执行入口
func (engine *Engine) Start(req *pb.EMReq) (rsp *pb.EMRsp, err error) {
	c := engine.pool.Get().(*Context)
	c.Req = req
	c.reset()

	rsp, err = engine.handle(c)

	engine.pool.Put(c)
	return
}

func (engine *Engine) handle(c *Context) (rsp *pb.EMRsp, err error) {
	rPath := c.Req.Cmd

	handlers := engine.getValue(rPath)
	if handlers != nil {
		c.handlers = handlers
		//按顺序执行中间件
		rsp, err = c.Next()
		return
	}
	return
}

//获取路由下的相关HandlersChain
func (engine *Engine) getValue(path string) (handlers HandlersChain) {
	handlers, ok := engine.tree[path]
	if !ok {
		return nil
	}
	return
}

func (engine *Engine) addRoute(path string, handlers HandlersChain) {
	engine.tree[path] = handlers
}

func (engine *Engine) Use(middleware ...HandlerFunc) {
	engine.Group.Use(middleware...)
}
