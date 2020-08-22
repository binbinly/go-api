package controller

import (
	"dj-api/app/common"
	pb "dj-api/proto"
)

//成功返回
func Suc(data []byte) *pb.EMRsp {
	if data == nil {
		data = []byte("")
	}
	return &pb.EMRsp{
		RspCode: common.Success,
		RspMsg:  "OK",
		RspData: data,
	}
}

//错误返回
func Err(msg string) *pb.EMRsp {
	return &pb.EMRsp{
		RspCode: common.Error,
		RspMsg:  msg,
	}
}
