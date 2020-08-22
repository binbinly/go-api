package agency

import (
	"dj-api/app/models"
)

type User struct {
	models.ModelUpdate

	UserId        int
	ParentId      int
	Channel       uint16
	Phone         string
	CurrentAmount int64
	DrawAmount    int64
	TotalAmount   int64
	RealName      string
	ChildCount    int
	Mobile        string
	QQ            int64
	Email         string
	Desc          string
	FreezeExplain string
	Status        int8
}

func (User) TableName() string {
	return "agency_user"
}
