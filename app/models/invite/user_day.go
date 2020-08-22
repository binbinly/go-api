package invite

import (
	"dj-api/app/models"
)

type UserDay struct {
	models.ModelOnlyUpdate

	UserId    int64
	ParentId  int64
	Date      string
	Flow      int64
	TotalFlow int64
}

func (UserDay) TableName() string {
	return models.TableInviteUserDay
}
