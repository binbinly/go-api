package invite

import (
	"dj-api/app/models"
)

type User struct {
	models.ModelUpdate

	UserId        int
	ParentId      int
	ChannelId     int
	Username      string
	Mobile        string
	CurrentAmount int64
	LastAmount    int64
	DrawAmount    int64
	TotalAmount   int64
	ChildCount    int
	TotalCount    int
	TotalFlow     int64
	TotalMyFlow   int64
	RegisterTime  int
}

func (User) TableName() string {
	return models.TableInviteUser
}
