package mq

import (
	"dj-api/app/mq/controller"
	"dj-api/dal/nsq/router"
)

func NsqRouter() *router.Engine {
	engine := router.NewEngine()

	engine.AddRoute("invite_user", controller.AddInviteUser)
	engine.AddRoute("invite_user_edit", controller.EditInviteUser)
	engine.AddRoute("user_order_flow", controller.OrderFlow)

	return engine
}
