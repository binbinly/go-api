package middleware

import (
	"dj-api/app/api/controller"
	"dj-api/app/api/middleware/auth/drivers"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *controller.User

		driver := drivers.NewTokenAuthDriver()
		if !driver.Check(c) { //检测token是否存在
			user = driver.Login(c)
			if user == nil {
				appG := controller.Gin{C: c}
				appG.ResponseError(controller.ErrorUserNotExist)
				return
			}
		} else {
			user = driver.User(c)
		}
		c.Set("channel_id", user.Channel)
		c.Set("user_id", user.UserId)
		c.Set("username", user.Username)
		c.Next()
	}
}
