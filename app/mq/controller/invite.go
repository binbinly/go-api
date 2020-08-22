package controller

import (
	"dj-api/app/logic"
	"dj-api/app/rpc/form"
	"dj-api/dal/nsq/router"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
)

/**
 * showdoc
 * @catalog nsq消息-agency
 * @title 添加代理
 * @description 添加代理
 * @url cmd: invite_user
 * @param user_id 必选 int 用户ID
 * @param parent_id 必选 int 上级ID
 * @param channel_id 必选 int 渠道ID
 * @param register_time 必选 int 注册时间
 * @param username 必选 string 昵称
 * @param mobile 必选 string 手机号
 * @remark 这里是备注信息
 * @number
 */
func AddInviteUser(c *router.Context) error {
	var userReq form.InviteUserReq

	err := json.Unmarshal(c.Req.Data, &userReq)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(userReq)
	if err != nil {
		return err
	}
	err = logic.AddInviteUser(userReq)
	if err != nil {
		return err
	}
	return nil
}

/**
 * showdoc
 * @catalog nsq消息-agency
 * @title 修改用户信息
 * @description 修改用户信息
 * @url cmd: invite_user_edit
 * @param user_id 必选 int 用户ID
 * @param username 必选 string 昵称
 * @param mobile 必选 string 手机号
 * @remark 这里是备注信息
 * @number
 */
func EditInviteUser(c *router.Context) error {
	var editUserReq form.InviteUserEdit

	err := json.Unmarshal(c.Req.Data, &editUserReq)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(editUserReq)
	if err != nil {
		return err
	}

	err = logic.EditInviteUser(editUserReq)
	if err != nil {
		return err
	}
	return nil
}
