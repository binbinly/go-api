package invite

import (
	"dj-api/app/models/invite"
	"dj-api/app/rpc/form"
	"dj-api/dal/db"
	"dj-api/tools"
	"github.com/jinzhu/gorm"
)

//添加一条记录
func AddUserDay(user form.OrderFlow) error {
	now := tools.DateFormat(tools.DateFormatString)
	var userDay invite.UserDay
	err := db.DB.Model(&invite.UserDay{}).Where("user_id = ? and date = ?", user.UserId, now).First(&userDay).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if userDay.ID > 0 { //已存在，修改
		err := db.DB.Model(&userDay).Updates(map[string]interface{}{
			"flow":       gorm.Expr("`flow` + ?", user.ValidFlow),
			"total_flow": gorm.Expr("`total_flow` + ?", user.Flow),
		}).Error
		if err != nil {
			return err
		}
	} else {
		data := &invite.UserDay{
			UserId:    user.UserId,
			ParentId:  user.SuperiorId,
			Flow:      user.ValidFlow,
			TotalFlow: user.Flow,
			Date:      now,
		}
		if err := db.DB.Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}
