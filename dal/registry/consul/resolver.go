package consul

import (
	"dj-api/tools/logger"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"net"
	"sync"
	"time"
)

const (
	defaultPort = "8500"
)

var (
	errMissingAddr = errors.New("consul resolver: missing address")

	errEndsWithColon = errors.New("consul resolver: missing port after port-separator colon")
)

//返回一个resolver.Builder的实例
func Init() {
	resolver.Register(NewBuilder())
}

//实现resolver.Builder的接口中的所有方法就是一个resolver.Builder
type consulBuilder struct {
}

type consulResolver struct {
	host                 string
	port                 string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	lastIndex            uint64
}

func NewBuilder() resolver.Builder {
	return &consulBuilder{}
}

//TODO 解析target, 拿到consul的ip和端口
//TODO 用consul的go api连接consul，查询服务结点信息，并且调用resolver.ClientConn的两个callback
func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	host, port, err := parseTarget(target.Authority, defaultPort)
	if err != nil {
		return nil, err
	}
	// IP address.
	if _, ok := formatIP(host); !ok {
		return nil, err
	}

	cr := &consulResolver{
		host:                 host,
		port:                 port,
		name:                 target.Endpoint,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		lastIndex:            0,
	}

	cr.wg.Add(1)
	go cr.watcher()
	return cr, nil
}

func (cr *consulResolver) watcher() {
	defer cr.wg.Done()

	for {
		services, metaInfo, err := consul.client.Health().Service(cr.name, "", true, &api.QueryOptions{WaitIndex: cr.lastIndex})
		if err != nil {
			logger.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}
		cr.lastIndex = metaInfo.LastIndex

		var newAddrs []resolver.Address

		for _, v := range services {
			newAddrs = append(newAddrs, resolver.Address{Addr: fmt.Sprintf("%v:%v", v.Service.Address, v.Service.Port)})
		}

		state := &resolver.State{
			Addresses: append(newAddrs),
		}
		cr.cc.UpdateState(*state)
	}
}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

//ResolverNow方法什么也不做，因为和consul保持了发布订阅的关系
//不需要像dns_resolver那个定时的去刷新
func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
}

//暂时先什么也不做吧
func (cr *consulResolver) Close() {
}

//解析
func parseTarget(target, defaultPort string) (host, port string, err error) {
	if target == "" {
		return "", "", errMissingAddr
	}
	if ip := net.ParseIP(target); ip != nil {
		// target is an IPv4 or IPv6(without brackets) address
		return target, defaultPort, nil
	}
	if host, port, err = net.SplitHostPort(target); err == nil {
		if port == "" {
			// If the port field is empty (target ends with colon), e.g. "[::1]:", this is an error.
			return "", "", errEndsWithColon
		}
		// target has port, i.e ipv4-host:port, [ipv6-host]:port, host-name:port
		if host == "" {
			// Keep consistent with net.Dial(): If the host is empty, as in ":80", the local system is assumed.
			host = "localhost"
		}
		return host, port, nil
	}
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil {
		// target doesn't have port
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
}

func formatIP(addr string) (addrIP string, ok bool) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", false
	}
	if ip.To4() != nil {
		return addr, true
	}
	return "[" + addr + "]", true
}
