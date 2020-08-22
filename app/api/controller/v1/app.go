package v1

import (
	"dj-api/app/api/controller"
	"dj-api/app/api/forms"
	"dj-api/app/rpc/service"
	"dj-api/app/rpc/service/php"
	"dj-api/dal/db"
	"dj-api/tools/logger"
	"github.com/gin-gonic/gin"
)

/**
 * showdoc
 * @catalog 基础接口
 * @title 获取渠道配置
 * @description 获取渠道配置
 * @method post
 * @url [url]/v1/config
 * @param channel 必选 int 渠道id
 * @return {"code": 200,"msg": "操作成功","data": {"cache_day": 2,"chat_live": "http://159.138.145.44:7005/chat.php?v=2&linkid=MTY0ZDg5ZGI2Zjg4NmY0MTZiZDJkMGE2MzU0MjA1Y2I_&ptn=","http_gateway": "http://120.25.239.229:8081","socket_gateway": "ws://120.25.239.229:8083/ws"}}
 * @return_param cache_day int 缓存天数
 * @return_param chat_live string 客服地址
 * @return_param http_gateway string http网关
 * @return_param socket_gateway string ws网关
 * @remark 这里是备注信息
 * @number
 */
func ChannelConfig(c *gin.Context) {
	var (
		appG = controller.Gin{C: c}
		form forms.ChannelConfig
	)

	isValid := controller.BindAndValid(c, &form)
	if !isValid {
		appG.ResponseError(controller.ErrorRequestParams)
		return
	}

	config, err := db.ChannelByCid(form.ChannelId)
	if err != nil {
		appG.ResponseError(controller.ErrorDataNotFound)
		return
	}
	appG.ResponseSuccess(config.AppConfig)
}

/**
 * showdoc
 * @catalog 基础接口
 * @title 版本升级
 * @description 版本升级
 * @method post
 * @url [url]/v1/upgrade
 * @param num 必选 string 当前版本号
 * @param platform 必选 int 平台
 * @return {"code": 200,"msg": "操作成功","data": {"version_number": "1","version_name": "压测版本1.3","desc": "adfadf","download_url": "","is_compel": 1,"ios_url": ""}}
 * @return_param version_number string 版本号
 * @return_param version_name string 版本名
 * @return_param desc string 描述
 * @return_param download_url string 下载地址
 * @return_param is_compel int 是否强制
 * @return_param ios_url string ios地址
 * @remark 这里是备注信息
 * @number
 */
func VersionUpgrade(c *gin.Context) {
	var (
		appG = controller.Gin{C: c}
		form forms.VersionUpgrade
	)

	isValid := controller.BindAndValid(c, &form)
	if !isValid {
		appG.ResponseError(controller.ErrorRequestParams)
		return
	}

	newVersion, err := db.Newest(form.Platform)
	if err != nil {
		appG.ResponseError(controller.ErrorDataNotFound)
		return
	}
	ret := controller.VersionCompare(form.Num, newVersion.VersionNumber)
	if ret == 0 {
		appG.ResponseError(controller.ErrorRequestParamsCheck)
		return
	} else if ret < 0 {
		appG.ResponseError(controller.ErrorNotNewVersion)
		return
	}
	appG.ResponseSuccess(newVersion)
}

/**
 * showdoc
 * @catalog 基础接口
 * @title 发送短信
 * @description 发送短信
 * @method post
 * @url [url]/v1/sms
 * @param type 必选 string 类型（可选值：login reg apply withdraw phone password）
 * @param mobile 必选 string 手机号
 * @return {"code": 200,"msg": "操作成功","data":{}}
 * @remark 这里是备注信息
 * @number
 */
func SendSms(c *gin.Context) {
	var (
		appG = controller.Gin{C: c}
		form forms.SendSms
	)

	isValid := controller.BindAndValid(c, &form)
	if !isValid {
		appG.ResponseError(controller.ErrorRequestParams)
		return
	}

	code, err := php.SmsSend(form.Mobile)
	if err != nil {
		logger.Warn(err)
		appG.ResponseError(controller.ERROR)
		return
	}
	//通知go服务
	request := service.NewBrgReq()
	request.SetCmd(service.CmdSms)
	_, err = request.Send(service.ParamSms{
		Phone: form.Mobile,
		Code:  code,
	})
	if err != nil {
		appG.ResponseError(controller.ErrorServiceError)
		return
	}
	appG.ResponseSuccessMsg("验证码：" + code)
}
