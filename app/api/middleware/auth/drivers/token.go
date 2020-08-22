package drivers

import (
	"dj-api/app/api/controller"
	"dj-api/app/common"
	"dj-api/app/rpc/service"
	"dj-api/dal/redis"
	"dj-api/tools/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type tokenAuthManager struct {
}

func NewTokenAuthDriver() *tokenAuthManager {
	return &tokenAuthManager{}
}

const (
	tokenExpire = 7200

	authToken = "auth_token"
)

func (t *tokenAuthManager) Check(c *gin.Context) bool {
	token := strings.TrimSpace(strings.Replace(c.Request.Header.Get("Authorization"), "Bearer", "", -1))
	if token == "" {
		return false
	}
	c.Set(authToken, token)
	if !redis.Client.IsExist(common.UserTokenPrefix + token) { //不存在
		return false
	}
	return true
}

func (t *tokenAuthManager) User(c *gin.Context) *controller.User {
	token := c.GetString(authToken)
	if token == "" {
		return nil
	}
	data, err := redis.Client.Get(common.UserTokenPrefix + token)
	if err != nil {
		logger.Warn(err)
		return nil
	}
	var user controller.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		logger.Warn(err)
		return nil
	}
	return &user
}

func (t *tokenAuthManager) Login(c *gin.Context) *controller.User {
	token := c.GetString(authToken)
	if token == "" {
		return nil
	}
	request := service.NewBrgReq()
	request.SetCmd(service.CmdCheckLogin)
	dataJson, err := request.Send(service.ParamCheckLogin{SessionId: token})
	if err != nil {
		logger.Info("check login err:%v", err)
		return nil
	}
	data := &service.RspCheckLogin{}
	err = json.Unmarshal(dataJson, &data)
	if err != nil {
		logger.Warn(err)
		return nil
	}
	user := &controller.User{
		UserId:   data.UserId,
		Channel:  data.RegisterChannel,
		Username: data.Phone,
	}

	str, err := json.Marshal(user)
	if err != nil {
		logger.Warn(err)
		return nil
	}
	redis.Client.Set(common.UserTokenPrefix+token, str, tokenExpire*time.Second)
	return user
}

func (t *tokenAuthManager) Logout(token string) bool {
	return redis.Client.Del(common.UserTokenPrefix + token)
}
