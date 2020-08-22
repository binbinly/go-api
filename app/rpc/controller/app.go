package controller

import (
	"dj-api/dal/db"
	"dj-api/dal/grpc/router"
	pb "dj-api/proto"
	"encoding/json"
	"fmt"
)

/**
 * showdoc
 * @catalog grpc接口
 * @title 获取团队信息
 * @description 获取团队信息
 * @url cmd: get_team
 * @return [{"channel_id":0,"team_id":3,"team_name":""},{"channel_id":0,"team_id":6,"team_name":""},{"channel_id":0,"team_id":3,"team_name":""},{"channel_id":0,"team_id":3,"team_name":""},{"channel_id":0,"team_id":3,"team_name":""},{"channel_id":0,"team_id":4,"team_name":""},{"channel_id":0,"team_id":3,"team_name":""},{"channel_id":0,"team_id":3,"team_name":""},{"channel_id":0,"team_id":3,"team_name":""},{"channel_id":0,"team_id":4,"team_name":""}]
 * @return_param 必选 int channel_id 渠道
 * @return_param 必选 int team_id 团队ID
 * @return_param 必选 string team_name 团队名
 * @remark 这里是备注信息
 * @number
 */
func GetTeam(c *router.Context) (pb *pb.EMRsp, err error) {
	var data []byte
	data, err = db.ChannelByTeam()
	if err != nil {
		return nil, err
	}
	return Suc(data), nil
}

type configReq struct {
	Name []string `json:"name" validate:"required"`
}

/**
 * showdoc
 * @catalog grpc接口
 * @title 获取配置
 * @description 获取配置 请求示例 ：[{"name":"cwt_rate"}]
 * @url cmd: get_config
 * @param name 必选 string 配置名
 * @return 返回结构：[{5}]
 * @return_param 必选 string value值
 * @remark 这里是备注信息
 * @number 数组结构
 */
func GetConfig(c *router.Context) (rsp *pb.EMRsp, err error) {
	var config configReq
	var data []byte

	err = json.Unmarshal(c.Req.ReqData, &config)
	if err != nil {
		err = fmt.Errorf("err:%v, req:%v", err, string(c.Req.ReqData))
		return
	}
	data, err = db.ConfigByNames(config.Name)
	if err != nil {
		return
	}
	return Suc(data), nil
}
