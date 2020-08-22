package invite

import (
	"database/sql/driver"
	"dj-api/app/models"
	"fmt"
	"time"
)

type Balance struct {
	UserId      int       `json:"user_id"`
	ParentId    int       `json:"parent_id"`
	Date        LocalTime `json:"date"`
	TotalFlow   int64     `json:"total_flow"`
	ChildFlow   int64     `json:"child_flow"`
	Rate        uint16    `json:"rate"`
	Amount      int64     `json:"amount"`
	TotalAmount int64     `json:"total_amount"`

	RegisterTime int `json:"register_time"`
}

type LocalTime struct {
	time.Time
}

func (Balance) TableName() string {
	return models.TableInviteBalance
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
