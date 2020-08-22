package agency

import (
	"dj-api/app/models/agency"
	"dj-api/dal/db"
)

//通过用户获取每期结算详情
func BalanceByUserId(userId int, time string) (*agency.Balance, error) {
	var balance agency.Balance
	err := db.DB.Model(&agency.Balance{}).Where("user_id = ? and `time` = ? ", userId, time).First(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}
