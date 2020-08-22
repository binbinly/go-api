package db

import (
	"dj-api/app/models"
	"dj-api/dal/redis"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	channelKey     = "ChannelModelConfig"
	channelTeamKey = "ChannelModelTeam"
)

var ChannelList map[int]models.SeoChannel

//渠道加载
func columnChannel() (err error) {
	var list []models.SeoChannel
	if !redis.Client.IsExist(channelKey) { //不存在，从数据库中获取
		err := DB.Model(&models.SeoChannel{}).Where("is_release = ? ", ReleaseYes).Select("id, team_id, app_config").Scan(&list).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if len(list) > 0 {
			ChannelList = make(map[int]models.SeoChannel, len(list))
			for _, v := range list {
				ChannelList[v.ID] = v
			}
			str, err := json.Marshal(ChannelList)
			if err != nil {
				return err
			}
			//加入缓存
			redis.Client.Set(channelKey, str, CacheTime)
		}
	} else {
		str, err := redis.Client.Get(channelKey)
		if err != nil {
			redis.Client.Del(channelKey)
			return err
		}
		err = json.Unmarshal([]byte(str), &ChannelList)
		if err != nil {
			redis.Client.Del(channelKey)
			return err
		}
	}
	return
}

//获取某个渠道配置
func ChannelByCid(cid int) (*models.SeoChannel, error) {
	err := columnChannel()
	if err != nil {
		return nil, err
	}
	data, ok := ChannelList[cid]
	if !ok {
		return nil, fmt.Errorf("channel:%v config not found", cid)
	}
	return &data, nil
}

//所有渠道对应团队
func ChannelByTeam() ([]byte, error) {
	var channelAll []models.ChannelTeam
	var data []byte
	var err error
	if !redis.Client.IsExist(channelTeamKey) {
		err = DB.Model(&models.SeoChannel{}).Select("seo_channel.id as channel_id, seo_channel.team_id, t.name as team_name").
			Where("is_release = ? ", ReleaseYes).
			Joins(" left join agency_team as t on t.id = seo_channel.team_id").Scan(&channelAll).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		data, err = json.Marshal(channelAll)
		if err != nil {
			return nil, err
		}
		redis.Client.Set(channelTeamKey, data, CacheTime)
	} else {
		str, err := redis.Client.Get(channelTeamKey)
		if err != nil {
			return nil, err
		}
		data = []byte(str)
	}
	return data, nil
}
