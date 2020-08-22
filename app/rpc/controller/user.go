package controller

import (
	"dj-api/app/common"
	"dj-api/dal/grpc/router"
	"dj-api/dal/redis"
	pb "dj-api/proto"
	"encoding/json"
	"fmt"
)

type logoutReq struct {
	SessionId string `json:"session_id"`
}

/**
 * showdoc
 * @catalog grpc接口
 * @title 用户注销
 * @description 用户注销
 * @url cmd: user_logout
 * @param session_id 必选 string 用户标识
 * @param flow 必选 int 流水
 * @remark 这里是备注信息
 * @number
 */
func Logout(c *router.Context) (pb *pb.EMRsp, err error) {
	var req logoutReq

	err = json.Unmarshal(c.Req.ReqData, &req)
	if err != nil {
		err = fmt.Errorf("err:%v, req:%v", err, string(c.Req.ReqData))
		return
	}
	if req.SessionId == "" {
		err = fmt.Errorf("session_id empty")
		return
	}

	redis.Client.Del(common.UserTokenPrefix + req.SessionId)
	return Suc(nil), nil
}
