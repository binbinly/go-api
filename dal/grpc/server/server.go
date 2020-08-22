package server

import (
	"context"
	"dj-api/app/config"
	"dj-api/app/rpc"
	"dj-api/dal/grpc/router"
	"dj-api/dal/registry"
	pb "dj-api/proto"
	"dj-api/tools"
	"dj-api/tools/logger"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

var ser *MyServer
var engine *router.Engine

type MyServer struct {
	host         string
	registry     registry.Registry
	service      *registry.Service
	gRpcServer   *grpc.Server
	healthServer *HealthImpl
}

func Start() error {
	var err error
	ser, err = RegisterService()
	if err != nil {
		return err
	}
	err = ser.ServiceStart()
	if err != nil {
		return err
	}
	registerRouter()
	return nil
}

func Stop() {
	if ser != nil {
		ser.registry.Unregister(nil, ser.service)
		ser.gRpcServer.GracefulStop()
	}
}

//注册路由
func registerRouter() {
	engine = rpc.GRpcRouter()
}

//服务注册
func RegisterService() (*MyServer, error) {
	ip, err := tools.GetLocalIP()
	if err != nil {
		return nil, err
	}

	registryInst, err := registry.InitRegistry(context.Background(), "consul",
		registry.WithAddr([]string{config.C.Registry.Host}),
	)
	if err != nil {
		return nil, err
	}

	service := &registry.Service{
		Name: config.C.Registry.ServiceName,
	}

	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   ip,
		Port: config.C.Registry.Port,
	})
	err = registryInst.Register(nil, service)
	if err != nil {
		return nil, err
	}

	var registrar = &MyServer{
		host:     ip,
		registry: registryInst,
		service:  service,
	}
	return registrar, nil
}

func (r *MyServer) ServiceStart() error {

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", r.host, config.C.Registry.Port))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	healthServer := &HealthImpl{}
	pb.RegisterEMServiceServer(s, &emServer{})
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	r.gRpcServer = s
	r.healthServer = healthServer

	go func() {
		if err := s.Serve(lis); err != nil {
			fmt.Printf("grpc server start err:%v\v", err)
		}
	}()
	return nil
}

type emServer struct{}

func (s *emServer) UnaryCall(ctx context.Context, req *pb.EMReq) (rsp *pb.EMRsp, err error) {
	rsp, err = engine.Start(req)
	if err != nil {
		logger.ErrorR(err)
	}
	return
}

func (s *emServer) ServerStreamCall(req *pb.EMReq, stream pb.EMService_ServerStreamCallServer) error {
	return nil
}

func (s *emServer) GameDataCall(req *pb.EMReq, stream pb.EMService_GameDataCallServer) error {
	return nil
}
