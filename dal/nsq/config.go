package nsq

import (
	"github.com/nsqio/go-nsq"
	"time"
)

func setting() *nsq.Config {
	c := nsq.NewConfig()
	c.MaxInFlight = 9
	c.LookupdPollInterval = 5 * time.Minute
	return c
}
