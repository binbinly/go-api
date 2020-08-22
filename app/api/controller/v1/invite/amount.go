package invite

import (
	"dj-api/app/api/controller"
	"dj-api/app/logic"
	"dj-api/app/rpc/service"
	"dj-api/dal/db"
	"dj-api/dal/db/invite"
	"dj-api/tools"
	"dj-api/tools/logger"
	"github.com/gin-gonic/gin"
	"time"
)

/**
 * showdoc
 * @catalog 邀请好友
 * @title 领取奖励
 * @description 领取奖励
 * @method post
 * @url [url]/v1/invite/draw
 * @return {"code": 200,"msg": "操作成功","data": {}}
 * @remark code = 200成功，其他皆失败
 * @number
 */
func Draw(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	appG := controller.Gin{C: c}

	user, err := invite.UserById(userId)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	if user == nil || user.ID == 0 {
		appG.ResponseError(controller.ErrorDataNotFound)
		return
	}
	//读取配置
	channelConfig, err := db.ChannelByCid(user.ChannelId)
	if err != nil {
		appG.ResponseError(controller.ErrorDataNotFound)
		return
	}

	config, err := invite.ColumnByTeamId(channelConfig.TeamId)
	if err != nil {
		appG.ResponseError(controller.ErrorDataNotFound)
		return
	}
	//当前佣金
	amount := user.CurrentAmount
	//金豆数量限额
	if amount < config.GoldLimit {
		appG.ResponseError(controller.ErrorInviteGoldLimit)
		return
	}
	//时间限制
	start, err := tools.FormatToTime(tools.DateFormat("2006-01-02 ") + config.StartTimeLimit)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	end, err := tools.FormatToTime(tools.DateFormat("2006-01-02 ") + config.EndTimeLimit)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}

	now := time.Now().Unix()
	if now < start || now > end {
		appG.ResponseError(controller.ErrorInviteTimeLimit)
		return
	}
	//开关限制
	if config.OpenLimit == 0 {
		appG.ResponseError(controller.ErrorInviteOpenLimit)
		return
	}

	request := service.NewBrgReq()
	request.SetCmd(service.CmdEditGold)
	_, err = request.Send(service.ParamEditGold{
		Event:    service.EditGoldTypeAdd,
		Gold:     amount,
		UserId:   userId,
		GoldType: 2,
	})
	if err != nil {
		logger.Warn(err)
		appG.ResponseError(controller.ErrorServiceError)
		return
	}
	if err := logic.UserDrawGold(user, amount); err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ErrorServer)
		return
	}
	appG.ResponseSuccessNil()
}

type LogItem struct {
	Time   string `json:"time"`
	Amount int64  `json:"amount"`
}

/**
 * showdoc
 * @catalog 邀请好友
 * @title 奖励记录
 * @description 奖励记录
 * @method post
 * @url [url]/v1/invite/log
 * @param p 非必选 int 第几页
 * @return {"code": 200,"msg": "操作成功","data": [{"time": "2020-05-30","amount": 10}]}
 * @return_param time string 时间
 * @return_param amount string 领取数额
 * @remark 这里是备注信息
 * @number
 */
func Log(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	appG := controller.Gin{C: c}

	logs, err := invite.LogById(userId, controller.GetPageOffset(c), controller.PageSize)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	data := make([]*LogItem, 0)
	for _, v := range logs {
		data = append(data, &LogItem{
			Time:   tools.DateUnixFormat(int64(v.CreatedAt), tools.DateTimeFormatStringSec),
			Amount: v.Amount,
		})
	}
	appG.ResponseSuccess(data)
}
