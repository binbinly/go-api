package agency

import (
	"dj-api/app/models"
)

type Apply struct {
	models.ModelUpdate

	UserId    int    `json:"user_id"`
	ChannelId int    `json:"channel_id"`
	RealName  string `json:"real_name"`
	Mobile    string `json:"mobile"`
	QQ        int    `json:"qq"`
	Email     string `json:"email"`
	Desc      string `json:"desc"`
	Status    int8   `json:"status"`
}

func (Apply) TableName() string {
	return "agency_apply"
}
