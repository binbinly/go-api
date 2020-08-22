package exception

import (
	"dj-api/app/api/controller"
	"dj-api/dal/grpc/router"
	pb "dj-api/proto"
	"dj-api/tools/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.ErrorR(fmt.Errorf("gin:%v", err))

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": controller.ErrorServer,
					"msg":  controller.MsgFlags[controller.ErrorServer],
				})
			}
		}()
		c.Next()
	}
}

func HandleRpcErrors() router.HandlerFunc {
	return func(c *router.Context) (pb *pb.EMRsp, err error) {
		defer func() {
			if e := recover(); e != nil {
				logger.ErrorR(fmt.Errorf("gRpc:%v", e))
				err = fmt.Errorf("panic:%v", e)
			}
		}()
		pb, err = c.Next()
		return
	}
}
