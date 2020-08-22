package nsq

import (
	"dj-api/app/config"
	"github.com/nsqio/go-nsq"
)

var producer *nsq.Producer

func Start() (err error) {
	producer, err = nsq.NewProducer(config.C.Nsq.ProdHost, setting())
	if err != nil {
		return err
	}
	return
}

func Stop() {
	if producer != nil {
		producer.Stop()
	}
}
