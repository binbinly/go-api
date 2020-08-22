package invite

import (
	"dj-api/app/api/controller"
	"dj-api/dal/db"
	"dj-api/dal/db/invite"
	"dj-api/tools/logger"
	"github.com/gin-gonic/gin"
)

const UserIdKey = "user_id"

type rspData struct {
	CurrentAmount int64 `json:"current_amount"`
	LastAmount    int64 `json:"last_amount"`
	ParentId      int   `json:"parent_id"`
	UserId        int   `json:"user_id"`
	ChildCount    int   `json:"child_count"`
	TotalAmount   int64 `json:"total_amount"`
	DrawAmount    int64 `json:"draw_amount"`
}

/**
 * showdoc
 * @catalog 邀请好友
 * @title 首页详情
 * @description 首页详情
 * @method get
 * @url [url]/v1/invite/home
 * @return {"code": 200,"msg": "操作成功","data": {"current_amount": 0,"parent_name": "","child_count": 0,"total_amount": 0,"draw_amount": 0}}
 * @return_param current_amount int 可领取奖励
 * @return_param last_amount int 上期奖励
 * @return_param parent_id int 上级id
 * @return_param user_id int 自己id
 * @return_param child_count int 直属好友
 * @return_param total_amount int 总奖励
 * @return_param draw_amount int 已领取奖励
 * @remark 这里是备注信息
 * @number
 */
func Home(c *gin.Context) {
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

	data := &rspData{
		CurrentAmount: user.CurrentAmount,
		LastAmount:    user.LastAmount,
		ChildCount:    user.ChildCount,
		ParentId:      user.ParentId,
		UserId:        userId,
		TotalAmount:   user.TotalAmount,
		DrawAmount:    user.DrawAmount,
	}
	appG.ResponseSuccess(data)
}

/**
 * showdoc
 * @catalog 邀请好友
 * @title 邀请配置
 * @description 邀请配置
 * @method get
 * @url [url]/v1/invite/config
 * @return {"code": 200,"msg": "操作成功","data": {"start_flow": 200,"end_flow": 300,"rate": 130,"name": ""}}
 * @return_param start_flow int 起始流水
 * @return_param end_flow int 结束流水
 * @return_param rate string 比例
 * @return_param name int 名称
 * @remark 这里是备注信息
 * @number
 */
func Config(c *gin.Context) {
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
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	appG.ResponseSuccess(config.Value)
}
