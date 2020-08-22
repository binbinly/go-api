package invite

import (
	"dj-api/app/models/invite"
	"dj-api/dal/db"
	"github.com/jinzhu/gorm"
)

func NewLog(userId int, amount int64, status int8) *invite.Log {
	return &invite.Log{
		UserId: userId,
		Amount: amount,
		Status: status,
	}
}

//根据用户获取领取记录
func LogById(userId, offset, limit int) (log []invite.Log, err error) {
	err = db.DB.Model(&invite.Log{}).Select("amount, created_at").
		Scopes(db.OffsetPage(offset, limit)).
		Where("user_id = ? and `status` = ? ", userId, db.StatusSuccess).Order(db.DefaultOrder).Find(&log).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}
