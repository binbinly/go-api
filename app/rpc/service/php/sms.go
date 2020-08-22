package php

import (
	"fmt"
)

const PathSms = "/sms/send"

//发送短信
func SmsSend(mobile string) (string, error) {
	data := map[string]interface{}{
		"mobile": mobile,
	}
	rsp, err := HttpGet(PathSms, data)
	if err != nil {
		return "", fmt.Errorf("sms send err:%v", err)
	}
	if code, ok := rsp["code"]; ok {
		return code.(string), nil
	}
	return "", fmt.Errorf("sms send err:%v, rsp data:%v", data, rsp)
}
