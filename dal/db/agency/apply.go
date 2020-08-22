package agency

import (
	"dj-api/app/api/forms"
	"dj-api/app/models/agency"
	"dj-api/dal/db"
	"html"
)

//添加一条申请记录
func AddApply(userId int, channelId int, form *forms.AddAgencyForm) error {
	apply := agency.Apply{
		UserId:    userId,
		ChannelId: channelId,
		RealName:  form.RealName,
		Mobile:    form.Mobile,
		QQ:        form.QQ,
		Email:     form.Email,
		Desc:      html.EscapeString(form.Desc),
	}
	if err := db.DB.Create(&apply).Error; err != nil {
		return err
	}
	return nil
}

//某用户是否存在
func ExistApply(userId int) bool {
	count := 0
	if err := db.DB.Model(&agency.Apply{}).Where("user_id = ? and status = ? ", userId, db.StatusInit).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

//获取最后一条申请记录
func ApplyById(userId int) (*agency.Apply, error) {
	var apply agency.Apply
	err := db.DB.Model(&agency.Apply{}).Where("user_id = ?", userId).Last(&apply).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}
