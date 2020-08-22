package invite

import (
	"dj-api/app/models"
)

type Log struct {
	models.ModelCreate

	UserId int   `json:"User_id"`
	Amount int64 `json:"amount"`
	Status int8  `json:"status"`
	Err    string
}

func (Log) TableName() string {
	return models.TableInviteLog
}
