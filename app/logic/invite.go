package logic

import (
	invitem "dj-api/app/models/invite"
	"dj-api/app/rpc/form"
	"dj-api/dal/db"
	"dj-api/dal/db/invite"
	"dj-api/dal/redis"
)

const QueueUserFlow = "queue:user_flow"

func AddInviteUser(user form.InviteUserReq) error {
	// 开始事务
	tx := db.DB.Begin()

	//添加一个用户
	err := invite.AddUser(user, tx)
	if err != nil {
		return err
	}
	err = invite.AddNode(user.UserId, user.ParentId, tx)
	if err != nil {
		return err
	}
	// 或提交事务
	tx.Commit()
	return nil
}

func EditInviteUser(user form.InviteUserEdit) error {
	err := invite.EditUser(user)
	if err != nil {
		return err
	}
	return nil
}

func AddInviteUserFlow(user form.OrderFlow) error {
	if err := invite.AddUserDay(user); err != nil {
		return err
	}
	//压入redis队列
	redis.Client.LPush(QueueUserFlow, user.UserId)
	return nil
}

//领取
func UserDrawGold(user *invitem.User, amount int64) error {
	log := invite.NewLog(user.UserId, amount, db.StatusInit)
	newAgency := map[string]int64{
		"CurrentAmount": 0,
		"DrawAmount":    user.DrawAmount + amount,
	}

	//开启事务
	tx := db.DB.Begin()
	err := tx.Model(&user).UpdateColumns(newAgency).Error
	if err != nil {
		log.Err = err.Error()
		log.Status = db.StatusError
		return err
	} else {
		log.Status = db.StatusSuccess
		err = tx.Model(&invitem.Log{}).Create(&log).Error
		if err != nil {
			tx.Rollback()
			return err
		} else {
			tx.Commit()
			return nil
		}
	}
}
