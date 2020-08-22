package service

import (
	"context"
	"dj-api/app/common"
	"dj-api/dal/grpc/meta"
	pb "dj-api/proto"
	"encoding/json"
	"fmt"
	"time"
)

const (
	CmdCheckLogin = "check_login"
	CmdEditGold   = "edit_gold"
	CmdSms        = "sms"

	EditGoldTypeAdd = "add" //加
	EditGoldTypeSub = "sub" //减
	EditGoldTypeSet = "set" //修改为指定值
)

type Request struct {
	ServiceName string      `json:"serviceName"`
	Cmd         string      `json:"cmd"`
	Data        interface{} `json:"data"`
}

func NewBrgReq() *Request {
	return &Request{
		ServiceName: serviceName,
	}
}

func (p *Request) SetCmd(cmd string) {
	p.Cmd = cmd
}

//gRpc远程调用
func (p *Request) Send(data interface{}) (rspData json.RawMessage, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = meta.InitRpcMeta(ctx, serviceName)

	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	rsp, err := cli.UnaryRequest(ctx, pb.EMReq{
		Cmd:     p.Cmd,
		ReqData: dataJson,
	})
	if err != nil {
		return nil, err
	}
	if rsp.RspCode == common.Success {
		if len(rsp.RspData) != 0 {
			rspData = rsp.RspData
		}
		return
	} else {
		err = fmt.Errorf("go service code: %d, msg:%s, reqData:%v", rsp.RspCode, rsp.RspMsg, data)
		return
	}
}
