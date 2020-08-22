package invite

import (
	"dj-api/app/api/controller"
	"dj-api/app/api/forms"
	"dj-api/dal/db/invite"
	"dj-api/tools"
	"dj-api/tools/logger"
	"github.com/gin-gonic/gin"
)

type balanceItem struct {
	UserId      int    `json:"user_id"`
	Flow        int64  `json:"flow"`
	Rate        uint16 `json:"rate"`
	Amount      int64  `json:"amount"`
	TotalAmount int64  `json:"total_amount"`
}

type dayBalanceRsp struct {
	List []balanceItem `json:"list"`
	Info balanceItem   `json:"info"`
}

/**
 * showdoc
 * @catalog 邀请好友
 * @title 每日记录
 * @description 每日记录
 * @method post
 * @url [url]/v1/invite/balance
 * @param p 非必选 int 第几页
 * @param date 非必选 string 日期
 * @param user_id 非必选 int 搜索用户id
 * @return {"code": 200,"msg": "操作成功","data": {"list": [],"info": {"user_id": 0,"flow": 0,"rate": 0,"amount": 0}}}
 * @return_param list int 下级列表
 * @return_param info int 自身信息
 * @return_param user_id int 用户ID
 * @return_param flow int 流水
 * @return_param rate int 比例
 * @return_param amount string 个人奖励
 * @return_param total_amount string 团队奖励
 * @remark 这里是备注信息
 * @number
 */
func Balance(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	var (
		appG = controller.Gin{C: c}
		form forms.InviteDetail
	)
	controller.BindAndValid(c, &form)
	if form.Date == "" {
		form.Date = tools.DateFormat(tools.DateFormatString)
	}
	my, err := invite.BalanceMyInfo(userId, form.Date)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	var r = &dayBalanceRsp{
		Info: balanceItem{
			UserId:      userId,
			Rate:        my.Rate,
			Flow:        my.TotalFlow,
			TotalAmount: my.TotalAmount,
			Amount:      my.Amount,
		},
		List: make([]balanceItem, 0),
	}
	if my == nil || my.UserId == 0 {
		appG.ResponseSuccess(r)
		return
	}

	balances, err := invite.BalanceChildByUserId(form.Date, userId, form.UserId, controller.GetPageOffset(c), controller.PageSize)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	for _, v := range balances {
		r.List = append(r.List, balanceItem{
			UserId:      v.UserId,
			Rate:        v.Rate,
			Flow:        v.TotalFlow,
			TotalAmount: v.TotalAmount,
			Amount:      my.Amount,
		})
	}
	appG.ResponseSuccess(r)
}

type balanceChildItem struct {
	UserId int    `json:"user_id"`
	Flow   int64  `json:"flow"`
	RTime  string `json:"r_time"`
}

/**
 * showdoc
 * @catalog 邀请好友
 * @title 流水详情
 * @description 流水详情
 * @method post
 * @url [url]/v1/invite/balance_detail
 * @param p 非必选 int 第几页
 * @param date 非必选 string 日期
 * @param user_id 非必选 int 搜索用户id
 * @return {"code": 200,"msg": "操作成功","data": [{"user_id": "15511","flow": 12222,"r_time": 130}]}
 * @return_param user_id int 用户id
 * @return_param r_time string 注册时间
 * @return_param flow int 流水
 * @remark 这里是备注信息
 * @number
 */
func BalanceDetail(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	var (
		appG = controller.Gin{C: c}
		form forms.InviteDetail
	)

	controller.BindAndValid(c, &form)
	if form.Date == "" {
		form.Date = tools.DateFormat(tools.DateFormatString)
	}

	tree, err := invite.TreeById(userId)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	if tree == nil || tree.ID == 0 {
		appG.ResponseError(controller.ErrorDataNotFound)
		return
	}

	balances, err := invite.BalanceChildAll(form.Date, tree.Root, tree.Left, tree.Right, form.UserId, controller.GetPageOffset(c), controller.PageSize)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	var data = make([]balanceChildItem, 0)
	for _, v := range balances {
		data = append(data, balanceChildItem{
			UserId: v.UserId,
			Flow:   v.ChildFlow,
			RTime:  tools.DateUnixFormat(int64(v.RegisterTime), tools.DateTimeFormatStringSec),
		})
	}
	appG.ResponseSuccess(data)
}
