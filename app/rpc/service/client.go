package service

import (
	"dj-api/dal/grpc/client"
)

const serviceName = "djBmBridge"

var cli *client.Client

func Init() (err error) {
	cli, err = client.NewClient(serviceName)
	return
}

func Close() (err error) {
	if cli != nil {
		return cli.Close()
	}
	return
}
