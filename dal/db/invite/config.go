package invite

import (
	"dj-api/app/models/invite"
	"dj-api/dal/db"
	"dj-api/dal/redis"
	"encoding/json"
	"fmt"
)

const (
	configKey = "InviteModelConfig"
)

var configList map[int]invite.Config

//配置加载
func columnConfig() (err error) {
	var list []invite.Config
	if !redis.Client.IsExist(configKey) { //不存在，从数据库中获取
		err = db.DB.Model(&invite.Config{}).Scan(&list).Error
		if err != nil {
			return
		}

		if len(list) > 0 {
			configList = make(map[int]invite.Config, len(list))
			for _, v := range list {
				configList[v.TeamId] = v
			}
			str, err := json.Marshal(configList)
			if err != nil {
				return err
			}
			redis.Client.Set(configKey, str, db.CacheTime)
		}
	} else {
		str, err := redis.Client.Get(configKey)
		if err != nil {
			redis.Client.Del(configKey)
			return err
		}
		err = json.Unmarshal([]byte(str), &configList)
		if err != nil {
			redis.Client.Del(configKey)
			return err
		}
	}
	return
}

//获取某个配置
func ColumnByTeamId(teamId int) (*invite.Config, error) {
	err := columnConfig()
	if err != nil {
		return nil, err
	}

	config, ok := configList[teamId]
	if !ok {
		return nil, fmt.Errorf("invite:%v config not found", teamId)
	}
	return &config, nil
}
