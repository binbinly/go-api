package invite

import (
	"database/sql/driver"
	"dj-api/app/models"
	"encoding/json"
)

type Config struct {
	models.ModelId

	TeamId         int    `json:"team_id"`
	Value          Values `json:"value"`
	GoldLimit      int64  `json:"gold_limit"`
	StartTimeLimit string `json:"start_time_limit"`
	EndTimeLimit   string `json:"end_time_limit"`
	OpenLimit      int    `json:"open_limit"`
}

type Values []*Value

type Value struct {
	StartFlow int    `json:"start_flow"`
	EndFlow   int    `json:"end_flow"`
	Rate      int    `json:"rate"`
	Name      string `json:"name"`
}

// Value 实现方法
func (p *Values) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *Values) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), p)
}

func (Config) TableName() string {
	return models.TableInviteConfig
}
