package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type RolePermissionRelService struct {
}

func (self *RolePermissionRelService) GetPermissionIDsByRoleIDs(roleIDs ...interface{}) ([]int, error) {
	permissionIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.RolePermissionRel{}).Where("role_id in (?)", roleIDs...).Pluck("distinct permission_id ",&permissionIDs).Error
	if err != nil {
		return nil, err
	}
	return permissionIDs, nil
}

func (self *RolePermissionRelService) AssignPermissions(roleID int, permissionIDs []int) (*[]int, error) {
	tx := common.SodaMngDB_WR.Begin()
	err := tx.Unscoped().Delete(permission.RolePermissionRel{}, "role_id = ?", roleID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range permissionIDs {
		err := tx.Create(&permission.RolePermissionRel{
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
