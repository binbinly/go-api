package v1

import (
	"dj-api/app/api/controller"
	"dj-api/dal/redis"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"github.com/unknwon/com"
	"time"
)

//测试使用
func AddUser(c *gin.Context) {
	appG := controller.Gin{C: c}

	userId := com.StrTo(c.Query("id")).MustInt()
	cid := com.StrTo(c.Query("c")).MustInt()
	key, _ := uuid.GenerateUUID()
	user := &controller.User{
		UserId:   userId,
		Channel:  cid,
		Username: "test",
	}
	userStr, err := json.Marshal(user)
	if err != nil {
		appG.ResponseError(controller.ErrorRequestParams)
		return
	}
	redis.Client.Set("token:"+key, userStr, time.Second*3600)
	appG.ResponseSuccess(key)
}
