package agency

import (
	"dj-api/app/models/agency"
	"dj-api/dal/db"
	"github.com/jinzhu/gorm"
)

//根据用户获取领取记录
func LogById(userId, offset, limit int) (list []agency.AmountDrawLog, err error) {
	err = db.DB.Model(&agency.AmountDrawLog{}).Select("amount, created_at").
		Where("user_id = ? and status = ? ", userId, db.StatusSuccess).Order(db.DefaultOrder).Scopes(db.OffsetPage(offset, limit)).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}
