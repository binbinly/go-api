package middleware

import (
	"dj-api/dal/grpc/router"
	pb "dj-api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Allow() bool
}

func Limit(l Limiter) router.HandlerFunc {
	return func(c *router.Context) (pb *pb.EMRsp, err error) {
		allow := l.Allow()
		if !allow {
			err = status.Error(codes.ResourceExhausted, "rate limited")
			return
		}
		pb, err = c.Next()
		return
	}
}
