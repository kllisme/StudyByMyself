package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/common"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type RolePermissionRelService struct {
}

func (self *RolePermissionRelService) GetPermissionIDsByRoleIDs(roleIDs ...interface{}) ([]int, error) {
	permissionIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&mngModel.RolePermissionRel{}).Where("role_id in (?)", roleIDs...).Pluck("distinct permission_id ",&permissionIDs).Error
	if err != nil {
		return nil, err
	}
	return permissionIDs, nil
}

func (self *RolePermissionRelService) AssignPermissions(roleID int, permissionIDs []int) (*[]int, error) {
	tx := common.SodaMngDB_WR.Begin()
	err := tx.Unscoped().Delete(mngModel.RolePermissionRel{}, "role_id = ?", roleID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range permissionIDs {
		err := tx.Create(&mngModel.RolePermissionRel{
			RoleID:       roleID,
			PermissionID: v,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &permissionIDs, nil
}
