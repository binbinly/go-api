package middleware

import (
	"dj-api/dal/grpc/router"
	pb "dj-api/proto"
	"dj-api/tools/logger"
)

//记录访问日志
func AccessLog(c *router.Context) (pb *pb.EMRsp, err error) {
	logger.Access(map[string]interface{}{
		"app":      "gRpc",
		"cmd":      c.Req.Cmd,
		"trace_id": c.Req.TraceID,
		"seq":      c.Req.Seq,
		"req":      string(c.Req.ReqData),
	})
	pb, err = c.Next()
	return
}
