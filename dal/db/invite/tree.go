package invite

import (
	"dj-api/app/models"
	"dj-api/app/models/invite"
	"dj-api/dal/db"
	"github.com/jinzhu/gorm"
)

//通过用户获取详情
func TreeById(userId int) (*invite.Tree, error) {
	tree := &invite.Tree{}
	err := db.DB.Model(&invite.Tree{}).Where("id = ? ", userId).First(&tree).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tree, nil
}

//添加树节点
func AddNode(userId, targetId int, tx *gorm.DB) error {
	if targetId == 0 {
		return addRoot(userId, tx)
	} else {
		return addChild(userId, targetId, tx)
	}
}

//添加子节点
func addChild(userId, targetId int, tx *gorm.DB) error {
	var target invite.Tree
	err := tx.Model(&invite.Tree{}).Where("id = ?", targetId).First(&target).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&invite.Tree{}).Where("`left` > ? and root = ? ", target.Left, target.Root).UpdateColumn("left", gorm.Expr("`left` + ?", 2)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&invite.Tree{}).Where("`right` > ? and root = ?", target.Left, target.Root).UpdateColumn("right", gorm.Expr("`right` + ?", 2)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tree := invite.Tree{
		ModelId: models.ModelId{ID: userId},
		Parent:  targetId,
		Level:   target.Level + 1,
		Left:    target.Left + 1,
		Right:   target.Left + 2,
		Root:    target.Root,
	}
	if err := tx.Create(&tree).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

//添加根节点
func addRoot(userId int, tx *gorm.DB) error {
	data := invite.Tree{
		ModelId: models.ModelId{ID: userId},
		Parent:  0,
		Level:   1,
		Left:    1,
		Right:   2,
		Root:    userId,
	}
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
