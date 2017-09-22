package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/common"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type PermissionMenuRelService struct {
}

func (self *PermissionMenuRelService) GetMenuIDsByPermissionIDs(permissionIDs ...interface{}) ([]int, error) {
	menuIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&mngModel.PermissionMenuRel{}).Where("permission_id in (?)", permissionIDs...).Pluck("distinct menu_id", &menuIDs).Error
	if err != nil {
		return nil, err
	}
	return menuIDs, nil
}

func (self *PermissionMenuRelService) AssignMenus(permissionID int, menuIDs []int) (*[]int, error) {
	tx := common.SodaMngDB_WR.Begin()
	err := tx.Unscoped().Delete(mngModel.PermissionMenuRel{}, "permission_id = ?", permissionID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range menuIDs {
		err := tx.Create(&mngModel.PermissionMenuRel{
			PermissionID: permissionID,
			MenuID:       v,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &menuIDs, nil
}
