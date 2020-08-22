package agency

import (
	"dj-api/app/api/controller"
	"dj-api/app/api/forms"
	"dj-api/app/logic"
	"dj-api/dal/db"
	"dj-api/dal/db/agency"
	"dj-api/dal/redis"
	"dj-api/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

const UserIdKey = "user_id"

//代理首页
func Home(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	appG := controller.Gin{C: c}

	agencyUser, err := agency.UserById(userId)
	if err != nil {
		appG.ResponseError(controller.ErrorUserNotExist)
		return
	}

	var detail = make(map[string]interface{}, 8)
	var amount = 0.0
	var activeNum uint16 = 0

	//获取推荐人数
	childCount := 0

	var data = make(map[string]interface{}, 8)
	data["status"] = agencyUser.Status
	data["child_user"] = childCount
	data["active_user"] = activeNum
	data["current_amount"] = agencyUser.CurrentAmount
	data["now_amount"] = map[string]interface{}{
		"amount": amount,
		"detail": detail,
	}
	data["total_amount"] = agencyUser.TotalAmount
	data["draw_amount"] = agencyUser.DrawAmount

	appG.ResponseSuccess(data)
}

//代理申请
func Apply(c *gin.Context) {
	userId := c.GetInt(UserIdKey)
	channelId := c.GetInt("channel_id")

	var (
		appG = controller.Gin{C: c}
		form forms.AddAgencyForm
	)

	isValid := controller.BindAndValid(c, &form)
	if !isValid {
		appG.ResponseError(controller.ErrorRequestParams)
		return
	}

	//验证验证码
	smsCode, _ := redis.Client.Get(fmt.Sprintf("sms_code:%s", form.Mobile))
	if smsCode != form.Code {
		appG.ResponseError(controller.ErrorVerifyCode)
		return
	}

	//xss过滤
	p := bluemonday.NewPolicy()
	form.Desc = p.Sanitize(form.Desc)

	ret, errCode := logic.CheckApply(userId)
	if !ret {
		appG.ResponseError(errCode)
		return
	}
	err := agency.AddApply(userId, channelId, &form)
	if err != nil {
		appG.ResponseError(controller.ErrorServerTimeout)
		return
	}
	appG.ResponseSuccess(map[string]string{})
}

//代理状态
func Status(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	appG := controller.Gin{C: c}

	var status int8 = 0 //已经是代理

	_, err := agency.UserById(userId)
	if err != nil { //还不是代理
		status = 1

		apply, err := agency.ApplyById(userId)
		if err != nil && apply != nil && apply.Status == db.StatusInit { //代理申请中
			status = 2
		}
	}
	appG.ResponseSuccess(map[string]int8{"status": status})
}

//某月结算详情
func Detail(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	var (
		appG = controller.Gin{C: c}
		form forms.AgencyDetail
	)

	isValid := controller.BindAndValid(c, &form)
	if !isValid {
		appG.ResponseError(controller.ErrorRequestParams)
		return
	}

	balance, err := agency.BalanceByUserId(userId, form.Time)
	if err != nil {
		appG.ResponseError(controller.ErrorDataNotFound)
		return
	}

	data := map[string]interface{}{
		"win":           balance.Win,
		"discounts":     balance.Discounts,
		"platform_fee":  balance.PlatformFee,
		"fee":           balance.Fee,
		"real_gold":     balance.RealWin,
		"rate":          fmt.Sprintf("%d", balance.Rate) + "%",
		"agency_amount": balance.AgencyAmount,
		"award":         balance.Award,
		"amount":        balance.Amount,
	}
	appG.ResponseSuccess(data)
}

//领取奖励
func Draw(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	appG := controller.Gin{C: c}

	agencyUser, _ := agency.UserById(userId)
	if agencyUser == nil {
		appG.ResponseError(controller.ErrorUserNotExist)
		return
	}
	if agencyUser.Status != db.StatusSuccess {
		appG.ResponseError(controller.ErrorUserFreeze)
		return
	}
	//当前佣金
	drawAmount := agencyUser.CurrentAmount
	if drawAmount == 0 {
		appG.ResponseError(controller.ErrorDataNotFound)
		return
	}
}

//领取奖励记录
func Log(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	appG := controller.Gin{C: c}

	logs, err := agency.LogById(userId, controller.GetPageOffset(c), controller.PageSize)
	if err != nil && len(logs) > 0 {
		var data []map[string]interface{}
		for _, v := range logs {
			data = append(data, map[string]interface{}{
				"time":   tools.DateUnixFormat(int64(v.CreatedAt), "2006-01-02"),
				"amount": v.Amount,
				"status": "已到账",
			})
		}
		appG.ResponseSuccess(data)
	} else {
		appG.ResponseSuccess(logs)
	}
}
