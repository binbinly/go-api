package router

//中间件组
type Group struct {
	engine *Engine
}

func (group *Group) AddRoute(absolutePath string, handler HandlerFunc) {
	//建立路由和相关中间件组的绑定
	group.engine.addRoute(absolutePath, handler)
}
