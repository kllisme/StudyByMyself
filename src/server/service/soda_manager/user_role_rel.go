package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/common"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type UserRoleRelService struct {
}

func (self *UserRoleRelService) GetRoleIDsByUserID(userID int) ([]int, error) {
	//TODO Considering User Group defined roles
	roleIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&mngModel.UserRoleRel{}).Where("user_id = ?", userID).Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}

func (self *UserRoleRelService) AssignRoles(userID int, roleIDs []int) (*[]int, error) {
	tx := common.SodaMngDB_WR.Begin()
	err := tx.Delete(mngModel.UserRoleRel{}, "user_id = ?", userID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range roleIDs {
		err := tx.Create(&mngModel.UserRoleRel{
			UserID: userID,
			RoleID: v,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &roleIDs, nil
}

//根据用户或角色ID删除关系表
func (self *UserRoleRelService) DeleteRel(userID int, roleID int) error {
	err := common.SodaMngDB_WR.Delete(mngModel.UserRoleRel{}, "user_id = ? or role_id = ?", userID, roleID).Error
	return err
}
