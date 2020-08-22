package logic

import (
	"dj-api/app/api/controller"
	"dj-api/dal/db/agency"
)

func CheckApply(userId int) (bool, uint16) {
	existUser := agency.ExistUser(userId)
	if !existUser {
		return false, controller.ErrorAgencyRepeatApply
	}
	existApply := agency.ExistApply(userId)
	if !existApply {
		return false, controller.ErrorApplyAuditing
	}
	return true, 0
}
