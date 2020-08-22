package router

//中间件组
type Group struct {
	//存储定义的中间件
	Handlers HandlersChain
	engine   *Engine
}

func (group *Group) Use(middleware ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middleware...)
}

func (group *Group) AddRoute(absolutePath string, handlers ...HandlerFunc) {
	handlers = group.combineHandlers(handlers)
	//建立路由和相关中间件组的绑定
	group.engine.addRoute(absolutePath, handlers)
}

//将定义的公用中间件和路由相关的中间件合并
func (group *Group) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}
