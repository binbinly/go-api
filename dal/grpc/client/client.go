package client

import (
	"context"
	"dj-api/app/config"
	"dj-api/dal/grpc/meta"
	"dj-api/dal/registry/consul"
	pb "dj-api/proto"
	"dj-api/tools/logger"
	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	c    pb.EMServiceClient
	conn *grpc.ClientConn
}

//创建一个grpc客户端
func NewClient(serviceName string) (*Client, error) {
	//consul resolver初始化
	consul.Init()

	target := "consul://" + config.C.Registry.Host + "/" + serviceName

	// Set up a connection to the server.
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithDisableServiceConfig())
	if err != nil {
		return nil, err
	}
	c := pb.NewEMServiceClient(conn)

	return &Client{c: c, conn: conn}, nil
}

func (cli *Client) UnaryRequest(ctx context.Context, req pb.EMReq) (rsp *pb.EMRsp, err error) {
	req.TraceID = logger.GetTraceId(ctx)
	rpcMeta := meta.GetRpcMeta(ctx)
	//熔断器
	err = hystrix.Do(rpcMeta.ServiceName, func() (err error) {
		rsp, err = cli.c.UnaryCall(ctx, &req)
		return err
	}, nil)
	return
}

func (cli *Client) ServerStreamRequest(ctx context.Context, req pb.EMReq) (pb.EMService_ServerStreamCallClient, error) {
	req.TraceID = logger.GetTraceId(ctx)
	return cli.c.ServerStreamCall(ctx, &req)
}

func (cli *Client) GameDataRequest(ctx context.Context, req pb.EMReq) (pb.EMService_GameDataCallClient, error) {
	req.TraceID = logger.GetTraceId(ctx)
	return cli.c.GameDataCall(ctx, &req)
}

func (cli *Client) Close() error {
	return cli.conn.Close()
}
