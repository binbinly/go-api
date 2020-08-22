package router

import (
	"sync"
)

type Engine struct {
	tree map[string]HandlerFunc // tree为了简化做成了map路由路径完全匹配
	Group
	pool sync.Pool // 正常情况存在大量的上下文切换，所以使用一个临时对象存储
}

func NewEngine() *Engine {
	engine := &Engine{
		tree: make(map[string]HandlerFunc),
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
func (engine *Engine) Start(req *Request) (err error) {
	c := engine.pool.Get().(*Context)
	c.Req = req

	handler := engine.getValue(c.Req.Cmd)
	if handler != nil {
		c.handler = handler
		//按顺序执行中间件
		err = c.Run()
	}

	engine.pool.Put(c)
	return
}

//获取路由下的相关HandlersChain
func (engine *Engine) getValue(path string) (handler HandlerFunc) {
	handler, ok := engine.tree[path]
	if !ok {
		return nil
	}
	return
}

func (engine *Engine) addRoute(path string, handler HandlerFunc) {
	engine.tree[path] = handler
}
