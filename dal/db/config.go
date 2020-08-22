package db

import (
	"dj-api/app/models"
	"dj-api/dal/redis"
	"encoding/json"
	"errors"
)

var configList map[string]interface{}

const (
	configKey = "ConfigModelList"

	TypeString = 1 //字符串文本
	TypeInt    = 2 //整形文本
	TypeJson   = 3 //json
)

//配置加载
func configLoad() (err error) {
	var list []models.SystemConfig
	if !redis.Client.IsExist(configKey) { //不存在，从数据库中获取
		err = DB.Model(&models.SystemConfig{}).Select("name, value").Scan(&list).Error
		if err != nil {
			return
		}
		if len(list) > 0 {
			configList = make(map[string]interface{})
			for _, v := range list {
				configList[v.Name] = v.Value
			}
			str, err := json.Marshal(configList)
			if err != nil {
				return err
			}
			//加入缓存
			redis.Client.Set(configKey, str, CacheTime)
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
func ConfigByName(name string) (interface{}, error) {
	err := configLoad()
	if err != nil {
		return "", err
	}
	value, ok := configList[name]
	if !ok {
		return "", errors.New("config not found to: " + name)
	}
	return value, nil
}

//获取多个配置返回给其他服务
func ConfigByNames(names []string) ([]byte, error) {
	err := configLoad()
	if err != nil {
		return nil, err
	}
	var values = make(map[string]interface{})
	for _, name := range names {
		value, ok := configList[name]
		if !ok {
			return nil, errors.New("config not found to: " + name)
		}
		values[name] = value
	}
	valuesJson, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	return valuesJson, nil
}
