package controller

const (
	SUCCESS uint16 = 200
	ERROR   uint16 = 500

	ErrorUserNotExist uint16 = 203
	ErrorUserFreeze   uint16 = 204
	ErrorVerifyCode   uint16 = 205

	ErrorRequestParams      uint16 = 400
	ErrorRequestParamsCheck uint16 = 401
	ErrorDataNotFound       uint16 = 404
	ErrorServer             uint16 = 501
	ErrorServerTimeout      uint16 = 502
	ErrorServiceError       uint16 = 505

	ErrorNotNewVersion        uint16 = 1000
	ErrorAgencyRepeatApply    uint16 = 1001
	ErrorApplyAuditing        uint16 = 1002
	ErrorAgencyAmountNotLimit uint16 = 1003
	ErrorAgencyNotExist       uint16 = 1004

	ErrorInviteGoldLimit uint16 = 1100
	ErrorInviteTimeLimit uint16 = 1101
	ErrorInviteOpenLimit uint16 = 1102
)

var MsgFlags = map[uint16]string{
	SUCCESS: "操作成功",
	ERROR:   "操作失败",

	ErrorUserNotExist: "请重新登录",
	ErrorUserFreeze:   "账号已被冻结，请联系客服",
	ErrorVerifyCode:   "验证码错误",

	ErrorRequestParams:      "请求参数不存在",
	ErrorRequestParamsCheck: "请求参数不合法",
	ErrorDataNotFound:       "信息不存在哦",
	ErrorServer:             "服务器错误",
	ErrorServerTimeout:      "服务器超时，清稍后再试",
	ErrorServiceError:       "服务器异常",

	ErrorNotNewVersion:        "已经是最新版本了哦",
	ErrorAgencyRepeatApply:    "不可以重复申请哦",
	ErrorApplyAuditing:        "申请正在处理中，请耐心等待哦",
	ErrorAgencyAmountNotLimit: "未达领取标准哦",
	ErrorAgencyNotExist:       "还不是代理哦",

	ErrorInviteGoldLimit: "未达到领取限额",
	ErrorInviteTimeLimit: "未在领取时间段哦",
	ErrorInviteOpenLimit: "请等待官方审核完成哦",
}

func GetMsg(code uint16, customMsg string) string {
	if customMsg != "" {
		return customMsg
	}
	if msg, ok := MsgFlags[code]; ok {
		return msg
	}
	return "未知错误"
}
