package agency

import (
	"dj-api/app/models"
)

type AmountDrawLog struct {
	models.ModelCreate

	UserId int   `json:"User_id"`
	Amount int64 `json:"amount"`
	Status int8  `json:"status"`
	Err    string
}
