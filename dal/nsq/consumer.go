package nsq

import (
	"dj-api/app/config"
	"dj-api/app/mq"
	"dj-api/dal/nsq/router"
	"dj-api/tools"
	"dj-api/tools/logger"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
)

const (
	TopicOrder   = "order_flow" //订单
	ChannelOrder = "invite"     //订单订阅channel
)

var agency *nsq.Consumer
var order *nsq.Consumer
var engine *router.Engine

//注册路由
func RegisterRouter() {
	engine = mq.NsqRouter()
}

//启动消费者
func Setup() (err error) {
	if agency, err = CreateConsumer(config.C.Nsq.Topic, config.C.Nsq.Channel); err != nil {
		return err
	}
	if order, err = CreateConsumer(TopicOrder, ChannelOrder); err != nil {
		return err
	}
	RegisterRouter()
	return nil
}

func StopC() {
	if agency != nil {
		agency.Stop()
	}
	if order != nil {
		order.Stop()
	}
}

//创建消费者
func CreateConsumer(topic, channel string) (*nsq.Consumer, error) {
	c, err := nsq.NewConsumer(topic, channel, setting())
	if err != nil {
		return nil, err
	}
	if !tools.IsDev() {
		f, err := os.Create(tools.GetRootDir() + "/logs/nsq_" + topic + ".log")
		if err != nil {
			logger.Panic("gin log file create fatal:%v", err)
		}
		c.SetLogger(log.New(f, "", log.Flags()), nsq.LogLevelInfo)
	}

	c.AddHandler(&Handler{})
	err = c.ConnectToNSQLookupds(config.C.Nsq.ConsumerHost)
	if err != nil {
		return nil, err
	}
	return c, nil
}
