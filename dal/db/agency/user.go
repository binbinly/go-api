package agency

import (
	"dj-api/app/models/agency"
	"dj-api/dal/db"
)

//某用户是否存在
func ExistUser(userId int) bool {
	count := 0
	if err := db.DB.Model(&agency.User{}).Where("user_id = ? ", userId).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

//通过用户获取详情
func UserById(userId int) (*agency.User, error) {
	var user agency.User
	err := db.DB.Model(&agency.User{}).Where("user_id = ? ", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
