package invite

import (
	"dj-api/app/models"
	"dj-api/app/models/invite"
	"dj-api/dal/db"
	"github.com/jinzhu/gorm"
)

//通过用户获取每期结算详情
func BalanceByUserId(userId, offset, limit int) ([]invite.Balance, error) {
	var balance []invite.Balance
	err := db.DB.Model(&invite.Balance{}).Scopes(db.OffsetPage(offset, limit)).Where("user_id = ?", userId).Order(db.DefaultOrder).Find(&balance).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return balance, nil
}

//获取详情
func BalanceMyInfo(userId int, date string) (*invite.Balance, error) {
	var balance invite.Balance
	err := db.DB.Model(&invite.Balance{}).Where("user_id = ? and date = ?", userId, date).First(&balance).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &balance, nil
}

//下级信息
func BalanceChildByUserId(date string, userId, sUserId, offset, limit int) (balance []invite.Balance, err error) {
	if sUserId > 0 {
		err = db.DB.Model(&invite.Balance{}).Where("date = ? and user_id = ? and parent_id = ?", date, sUserId, userId).Find(&balance).Error
	} else {
		err = db.DB.Model(&invite.Balance{}).Where("date = ? and parent_id = ?", date, userId).Scopes(db.OffsetPage(offset, limit)).Order("id asc").Find(&balance).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return balance, nil
}

//下级所有详情
func BalanceChildAll(date string, root, left, right, sUserId, offset, limit int) (balance []invite.Balance, err error) {
	if sUserId > 0 {
		err = db.DB.Model(&invite.Balance{}).Where("u.user_id = ?", sUserId).
			Scopes(whereTree(date, root, left, right), joinUser).Find(&balance).Error
	} else {
		err = db.DB.Model(&invite.Balance{}).
			Where(models.TableInviteBalance+".child_flow > 0 ").
			Scopes(whereTree(date, root, left, right), db.OffsetPage(offset, limit), joinUser).Order("id asc").Find(&balance).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return balance, nil
}

func joinUser(db *gorm.DB) *gorm.DB {
	return db.Joins(" left join " + models.TableInviteUser + " as u on u.user_id = " + models.TableInviteBalance + ".user_id")
}

func whereTree(date string, root, left, right int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(models.TableInviteBalance+".*, u.register_time").
			Where("date = ? and "+models.TableInviteBalance+".user_id in ( SELECT id FROM "+models.TableInviteTree+" WHERE `root` = ? and `left` BETWEEN ? and ? ) ", date, root, left, right)
	}
}
