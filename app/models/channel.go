package models

import (
	"database/sql/driver"
	"encoding/json"
)

type SeoChannel struct {
	ModelId
	TeamId    int       `json:"team_id"`
	AppConfig appConfig `json:"app_config"`
}

type AppConfig struct {
	CacheDay      int    `json:"cache_day"`
	HttpGateway   string `json:"http_gateway"`
	SocketGateway string `json:"socket_gateway"`
}

type ChannelTeam struct {
	ChannelId int    `json:"channel_id"`
	TeamId    int    `json:"team_id"`
	TeamName  string `json:"team_name"`
}

type appConfig AppConfig

// Value 实现方法
func (p *appConfig) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *appConfig) Scan(input interface{}) error {
	var inputB = input.([]byte)
	if len(inputB) > 0 {
		return json.Unmarshal(inputB, p)
	}
	return nil
}
