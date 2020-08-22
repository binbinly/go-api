package registry

//服务抽象
type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes"`
}

//服务节点抽象
type Node struct {
	Weight int    `json:"weight"` //权重
	Port   int    `json:"port"`   //端口
	ID     string `json:"id"`
	IP     string `json:"ip"`
}
