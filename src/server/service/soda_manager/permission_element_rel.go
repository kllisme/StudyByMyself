package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/common"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type PermissionElementRelService struct {
}

func (self *PermissionElementRelService) GetElementIDsByPermissionIDs(permissionIDs ...interface{}) ([]int, error) {
	elementIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&mngModel.PermissionElementRel{}).Where("permission_id in (?)", permissionIDs...).Pluck("distinct element_id", &elementIDs).Error
	if err != nil {
		return nil, err
	}
	return elementIDs, nil
}

func (self *PermissionElementRelService) AssignElements(permissionID int, elementIDs []int) (*[]int, error) {
	tx := common.SodaMngDB_WR.Begin()
	err := tx.Unscoped().Delete(mngModel.PermissionElementRel{}, "permission_id = ?", permissionID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range elementIDs {
		err := tx.Create(&mngModel.PermissionElementRel{
			PermissionID: permissionID,
			ElementID:    v,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &elementIDs, nil
}
