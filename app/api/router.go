package api

import (
	"dj-api/app/api/controller"
	"dj-api/app/api/controller/v1"
	"dj-api/app/api/controller/v1/agency"
	"dj-api/app/api/controller/v1/invite"
	"dj-api/app/api/middleware"
	"dj-api/exception"
	"dj-api/tools"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func InitRouter() *gin.Engine {

	router := gin.New()

	if tools.IsDev() {
		pprof.Register(router) // 性能分析工具
	}

	router.Use(gin.Logger())
	router.Use(middleware.Cors())        //跨域
	router.Use(exception.HandleErrors()) // 错误处理

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该操作",
		})
	})

	registerApiRouter(router)
	return router
}

func registerApiRouter(router *gin.Engine) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", tools.Mobile)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("v_num", tools.VersionNum)
	}

	router.POST("/notify/slack/MBuIh3DcCtpRILPjqRuQduwC27JmSl63", controller.Notify)
	api := router.Group("/v1")

	api.POST("/config", v1.ChannelConfig)
	api.POST("/upgrade", v1.VersionUpgrade)
	api.POST("/sms", v1.SendSms)
	api.GET("/test_user", v1.AddUser)

	apiAgency := api.Group("/agency")
	apiAgency.Use(middleware.Auth())
	{
		apiAgency.POST("/apply", agency.Apply)
		apiAgency.GET("/home", agency.Home)
		apiAgency.GET("/status", agency.Status)
		apiAgency.POST("/detail", agency.Detail)
		apiAgency.GET("/log", agency.Log)
		apiAgency.POST("/draw", agency.Draw)
	}

	apiInvite := api.Group("/invite")
	apiInvite.Use(middleware.Auth())
	{
		apiInvite.GET("/home", invite.Home)
		apiInvite.GET("/config", invite.Config)
		apiInvite.POST("/detail", invite.Detail)
		apiInvite.POST("/draw", invite.Draw)
		apiInvite.POST("/log", invite.Log)
		apiInvite.POST("/balance", invite.Balance)
		apiInvite.POST("/balance_detail", invite.BalanceDetail)
	}
}
