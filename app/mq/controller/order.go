package controller

import (
	"dj-api/app/logic"
	"dj-api/app/rpc/form"
	"dj-api/dal/nsq/router"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
)

//用户产生订单
func OrderFlow(c *router.Context) error {
	var userReq form.OrderFlow

	err := json.Unmarshal(c.Req.Data, &userReq)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(userReq)
	if err != nil {
		return err
	}
	err = logic.AddInviteUserFlow(userReq)
	if err != nil {
		return err
	}
	return nil
}
