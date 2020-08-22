package nsq

import (
	"context"
	"dj-api/dal/nsq/router"
	"dj-api/tools/logger"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type msg struct {
	TraceID  string          `json:"trace_id"`       // 消息跟踪ID
	ExpireAt int64           `json:"expire_at"`      // 过期时间
	UserIDs  []int64         `json:"user_ids"`       // 用户ID，没有userID的消息不需要关心此字段
	Cmd      string          `json:"cmd"`            // 请求命令
	Seq      int32           `json:"seq"`            // 消息序号
	Data     json.RawMessage `json:"data,omitempty"` // 数据 json
}

func PushMsg(topic, cmd string, data json.RawMessage) error {
	msg := &msg{
		TraceID:  logger.GetTraceId(context.Background()),
		ExpireAt: 0,
		UserIDs:  nil,
		Cmd:      cmd,
		Seq:      0,
		Data:     data,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = producer.Publish(topic, msgJson)
	if err != nil {
		return err
	}
	return nil
}

type Handler struct{}

func (h *Handler) HandleMessage(m *nsq.Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.ErrorR(fmt.Errorf("nsq:%v", e))
			err = fmt.Errorf("panic:%v", e)
		}
	}()
	var msgContent msg

	err = json.Unmarshal(m.Body, &msgContent)
	if err != nil {
		logger.Info("nsq un marshal err :%v, body:%v", err, string(m.Body))
		return nil
	}
	err = engine.Start(&router.Request{
		Cmd:  msgContent.Cmd,
		Data: msgContent.Data,
	})
	if err != nil {
		logger.ErrorR(fmt.Errorf("nsq err :%v, body:%v", err, string(m.Body)))
	}
	return err
}
