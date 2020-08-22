package rpc

import (
	"dj-api/app/config"
	"dj-api/app/rpc/controller"
	"dj-api/app/rpc/middleware"
	"dj-api/dal/grpc/router"
	"dj-api/exception"
	"golang.org/x/time/rate"
)

func GRpcRouter() *router.Engine {
	engine := router.NewEngine()

	engine.Use(exception.HandleRpcErrors())
	engine.Use(middleware.AccessLog)
	//初始化限流器
	if config.C.Limit.Enable {
		limit := rate.NewLimiter(rate.Limit(config.C.Limit.Qps),
			config.C.Limit.Qps)
		engine.Use(middleware.Limit(limit))
	}

	engine.AddRoute("user_logout", controller.Logout)
	engine.AddRoute("get_config", controller.GetConfig)
	engine.AddRoute("get_team", controller.GetTeam)
	return engine
}
