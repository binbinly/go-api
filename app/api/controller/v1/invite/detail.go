package invite

import (
	"dj-api/app/api/controller"
	"dj-api/app/api/forms"
	"dj-api/dal/db/invite"
	"dj-api/tools"
	"dj-api/tools/logger"
	"github.com/gin-gonic/gin"
)

type rspDataDetail struct {
	UserId       int    `json:"user_id"`
	RegisterTime string `json:"register_time"`
	MyFlow       int64  `json:"my_flow"`
}

/**
 * showdoc
 * @catalog 邀请好友
 * @title 用户详情
 * @description 用户性情
 * @method post
 * @url [url]/v1/invite/detail
 * @param user_id 非必选 string 昵称
 * @param p 非必选 int 第几页
 * @return {"code": 200,"msg": "操作成功","data": [{"user_id": "15516","register_time": "1970-01-01","my_flow": 0}]}
 * @return_param user_id int 昵称
 * @return_param register_time string 注册时间
 * @return_param my_flow int 个人流水
 * @remark 这里是备注信息
 * @number
 */
func Detail(c *gin.Context) {
	userId := c.GetInt(UserIdKey)

	var (
		appG = controller.Gin{C: c}
		form forms.InviteDetail
	)

	controller.BindAndValid(c, &form)
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

	list, err := invite.ChildAllUserList(form.UserId, tree.Root, tree.Left, tree.Right, controller.GetPageOffset(c), controller.PageSize)
	if err != nil {
		logger.Error(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	data := make([]rspDataDetail, 0)
	for _, v := range list {
		var detail rspDataDetail
		detail.UserId = v.UserId
		detail.MyFlow = v.TotalMyFlow
		detail.RegisterTime = tools.DateUnixFormat(int64(v.RegisterTime), tools.DateFormatString)
		data = append(data, detail)
	}
	appG.ResponseSuccess(data)
}
