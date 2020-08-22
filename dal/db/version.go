package db

import (
	"dj-api/app/models"
	"dj-api/dal/redis"
	"encoding/json"
	"fmt"
)

const (
	TypeMain    = 1 //大版本类型包
	VersionInfo = "AppVersionModelInfo"
)

//获取最新版本
func Newest(platform int) (*models.AppVersion, error) {
	var newVersion models.AppVersion
	var cacheKey = fmt.Sprintf("%s_%d", VersionInfo, platform)

	if !redis.Client.IsExist(cacheKey) { //不存在，从数据库中获取
		err := DB.Model(&models.AppVersion{}).Where("is_release = ? and type = ? and platform = ? ", ReleaseYes, TypeMain, platform).
			Select("version_number, version_name, download_url, `desc`, is_compel, ios_url").Last(&newVersion).Error
		if err != nil {
			return nil, err
		}
		str, err := json.Marshal(newVersion)
		if err != nil {
			return nil, err
		}
		redis.Client.Set(cacheKey, str, CacheTime)
	} else {
		str, err := redis.Client.Get(cacheKey)
		if err != nil {
			redis.Client.Del(cacheKey)
			return nil, err
		}
		err = json.Unmarshal([]byte(str), &newVersion)
		if err != nil {
			redis.Client.Del(cacheKey)
			return nil, err
		}
	}
	return &newVersion, nil
}
