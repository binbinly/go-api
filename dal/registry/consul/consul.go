package consul

import (
	"context"
	"dj-api/dal/grpc/server"
	"dj-api/dal/registry"
	"dj-api/tools/logger"
	"fmt"
	"github.com/hashicorp/consul/api"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxServiceNum          = 8                //服务最大数量
	MaxSyncServiceInterval = time.Second * 10 //健康检查间隔
	Deregister             = time.Minute      //服务自动注销时间
)

//consul 注册插件
type Registry struct {
	options   *registry.Options
	client    *api.Client
	serviceCh chan *registry.Service

	value              atomic.Value
	lock               sync.Mutex
	registryServiceMap map[string]*RegisterService
}

//所有服务存储
type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

type RegisterService struct {
	registered bool
	id         string
	service    *registry.Service
}

var consul = &Registry{
	serviceCh:          make(chan *registry.Service, MaxServiceNum),
	registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
}

func init() {
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	consul.value.Store(allServiceInfo)
	registry.RegisterPlugin(consul)
	go consul.run()
}

//插件的名字
func (e *Registry) Name() string {
	return "consul"
}

//初始化
func (e *Registry) Init(ctx context.Context, opts ...registry.Option) (err error) {

	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}

	config := api.DefaultConfig()
	config.Address = e.options.Addr[0]
	e.client, err = api.NewClient(config)

	if err != nil {
		return err
	}
	return
}

//服务注册
func (e *Registry) Register(ctx context.Context, service *registry.Service) (err error) {

	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}
	return
}

//服务反注册
func (e *Registry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	for _, v := range e.registryServiceMap {
		err = e.client.Agent().ServiceDeregister(v.id)
	}
	return
}

//服务发现
func (e *Registry) Find(ctx context.Context, name string) (service *registry.Service, err error) {

	//一般情况下，都会从缓存中读取
	service, ok := e.getServiceFromCache(name)
	if ok {
		return
	}

	//如果缓存中没有这个service，则从consul中读取
	e.lock.Lock()
	defer e.lock.Unlock()
	//先检测，是否已经从consul中加载成功了
	service, ok = e.getServiceFromCache(name)
	if ok {
		return
	}

	services, _, err := e.client.Health().Service(name, "", true, &api.QueryOptions{})
	if err != nil {
		return
	}

	service = &registry.Service{
		Name: name,
	}

	for _, v := range services {
		node := registry.Node{
			Weight: v.Service.Weights.Passing,
			Port:   v.Service.Port,
			ID:     v.Service.ID,
			IP:     v.Service.Address,
		}
		service.Nodes = append(service.Nodes, &node)
	}

	allServiceInfoOld := e.value.Load().(*AllServiceInfo)
	var allServiceInfoNew = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	for key, val := range allServiceInfoOld.serviceMap {
		allServiceInfoNew.serviceMap[key] = val
	}

	allServiceInfoNew.serviceMap[name] = service
	e.value.Store(allServiceInfoNew)
	return
}

func (e *Registry) run() {
	ticker := time.NewTicker(MaxSyncServiceInterval)
	for {
		select {
		case service := <-e.serviceCh:
			registryService, ok := e.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes {
					registryService.service.Nodes = append(registryService.service.Nodes, node)
				}
				registryService.registered = false
				break
			}
			registryService = &RegisterService{
				service: service,
			}

			err := e.registerService(registryService)
			if err != nil {
				logger.Error(err)
			}

			e.registryServiceMap[service.Name] = registryService
		case <-ticker.C:
			if server.GRpcHealUpdateAt > 0 && time.Now().Unix()-server.GRpcHealUpdateAt > 20 {
				for _, s := range e.registryServiceMap {
					err := e.registerService(s)
					if err != nil {
						logger.Error(err)
					}
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

//注册进入consul
func (e *Registry) registerService(registryService *RegisterService) (err error) {
	registryService.id = fmt.Sprintf("%s-%d", registryService.service.Name, time.Now().Unix())
	for _, node := range registryService.service.Nodes {
		reg := &api.AgentServiceRegistration{
			ID:      fmt.Sprintf("%v-%v-%v", registryService.service.Name, node.IP, node.Port),
			Name:    registryService.service.Name,
			Tags:    []string{registryService.service.Name},
			Port:    node.Port,
			Address: node.IP,
			Check: &api.AgentServiceCheck{ // 健康检查
				Interval:                       MaxSyncServiceInterval.String(),                                           // 健康检查间隔
				GRPC:                           fmt.Sprintf("%v:%v/%v", node.IP, node.Port, registryService.service.Name), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
				DeregisterCriticalServiceAfter: Deregister.String(),                                                       // 注销时间，相当于过期时间
			},
		}
		if err = e.client.Agent().ServiceRegister(reg); err != nil {
			err = fmt.Errorf("consul service register fail: %s", err.Error())
			continue
		}
		registryService.registered = true
	}

	return
}

func (e *Registry) getServiceFromCache(name string) (service *registry.Service, ok bool) {
	allServiceInfo := e.value.Load().(*AllServiceInfo)
	//一般情况下，都会从缓存中读取
	service, ok = allServiceInfo.serviceMap[name]
	return
}

//从consul更新服务
func (e *Registry) syncServiceFromConsul() {

	var allServiceInfoNew = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	allServiceInfo := e.value.Load().(*AllServiceInfo)

	//对于缓存的每一个服务，都需要从consul中进行更新
	for _, service := range allServiceInfo.serviceMap {
		services, _, err := e.client.Health().Service(service.Name, "", true, &api.QueryOptions{})
		if err != nil {
			allServiceInfoNew.serviceMap[service.Name] = service
			continue
		}

		serviceNew := &registry.Service{
			Name: service.Name,
		}

		for _, v := range services {
			node := registry.Node{
				Weight: v.Service.Weights.Passing,
				Port:   v.Service.Port,
				ID:     v.Service.ID,
				IP:     v.Service.Address,
			}
			serviceNew.Nodes = append(serviceNew.Nodes, &node)
		}
		allServiceInfoNew.serviceMap[serviceNew.Name] = serviceNew
	}

	e.value.Store(allServiceInfoNew)
}
