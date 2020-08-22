package invite

import (
	"dj-api/app/models/invite"
	"dj-api/app/rpc/form"
	"dj-api/dal/db"
	"github.com/jinzhu/gorm"
)

//用户是否存在
func ExistUser(userId int) bool {
	count := 0
	if err := db.DB.Model(&invite.User{}).Where("user_id = ? ", userId).Count(&count).Error; err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

//通过用户获取详情
func UserById(userId int) (*invite.User, error) {
	var user invite.User
	err := db.DB.Model(&invite.User{}).Where("user_id = ? ", userId).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &user, nil
}

//所有子级列表
func ChildAllUserList(sUserId, root, left, right, offset, limit int) (list []invite.User, err error) {
	if sUserId != 0 {
		err = db.DB.Model(&invite.User{}).Select("user_id, register_time, total_my_flow").Where("user_id = ? ", sUserId).Find(&list).Error
	} else {
		err = db.DB.Model(&invite.User{}).Select("user_id, register_time, total_my_flow").
			Where("total_my_flow > 0 ").
			Where("user_id in ( SELECT id FROM invite_tree WHERE `root` = ? and `left` BETWEEN ? and ? ) ", root, left, right).
			Order("id asc").Scopes(db.OffsetPage(offset, limit)).Find(&list).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

//直属（包括自己）
func ChildUserList(sUserId, userId, offset, limit int) (list []invite.User, err error) {
	if sUserId != 0 {
		err = db.DB.Model(&invite.User{}).Select("user_id, register_time, child_count, total_flow").
			Where("user_id = ? ", sUserId).
			Order("id asc").Offset(offset).Limit(limit).Find(&list).Error
	} else {
		err = db.DB.Model(&invite.User{}).Select("user_id, register_time, child_count, total_flow").
			Where("parent_id = ? or user_id = ? ", userId, userId).
			Order("id asc").Scopes(db.OffsetPage(offset, limit)).Find(&list).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

//添加用户
func AddUser(user form.InviteUserReq, tx *gorm.DB) error {
	data := invite.User{
		UserId:       user.UserId,
		ChannelId:    user.ChannelId,
		Username:     user.Username,
		ParentId:     user.ParentId,
		RegisterTime: user.RegisterTime,
		Mobile:       user.Mobile,
	}
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return err
	}
	if user.ParentId > 0 {
		var parentUser invite.User
		err := tx.Model(&invite.User{}).Where("user_id = ? ", user.ParentId).First(&parentUser).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Model(&parentUser).UpdateColumn("child_count", gorm.Expr("child_count + ?", 1)).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

//修改用户
func EditUser(user form.InviteUserEdit) error {
	inviteUser, err := UserById(user.UserId)
	if err != nil {
		return err
	}
	if user.Username != "" {
		inviteUser.Username = user.Username
	}
	if user.Mobile != "" {
		inviteUser.Mobile = user.Mobile
	}
	err = db.DB.Save(inviteUser).Error
	if err != nil {
		return err
	}
	return nil
}
