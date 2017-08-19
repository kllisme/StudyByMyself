package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type PermissionElementRelService struct {

}

func (self *PermissionElementRelService) GetElementIDsByPermissionIDs(permissionIDs ...interface{}) ([]int, error) {
	elementIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.PermissionElementRel{}).Where("permission_id in (?)", permissionIDs...).Pluck("element_id", &elementIDs).Error
	if err != nil {
		return nil, err
	}
	return elementIDs, nil
}

func (self *PermissionElementRelService) AssignElements(permissionID int, elementIDs []int) (*[]int, error) {
	tx := common.SodaMngDB_R.Begin()
	err := tx.Unscoped().Delete(permission.PermissionElementRel{}, "permission_id = ?", permissionID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range elementIDs {
		err := tx.Create(&permission.PermissionElementRel{
			PermissionID:permissionID,
			ElementID:v,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &elementIDs, nil
}
