package router

import (
	"encoding/json"
)

type HandlerFunc func(*Context) error

type Request struct {
	Cmd  string
	Data json.RawMessage
}

//定义的上下文
type Context struct {
	Req     *Request
	handler HandlerFunc
}

//调用
func (c *Context) Run() (err error) {
	err = c.handler(c)
	return
}
