package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type PermissionActionRelService struct {

}

func (self *PermissionActionRelService) GetActionIDsByPermissionIDs(permissionIDs ...interface{}) ([]int, error) {
	actionIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.PermissionActionRel{}).Where("permission_id in (?)", permissionIDs...).Pluck("distinct action_id",&actionIDs).Error
	if err != nil {
		return nil, err
	}
	return actionIDs, nil
}

func (self *PermissionActionRelService) AssignActions(permissionID int, actionIDs []int) (*[]int, error) {
	tx := common.SodaMngDB_R.Begin()
	err := tx.Unscoped().Delete(permission.PermissionActionRel{}, "permission_id = ?", permissionID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range actionIDs {
		err := tx.Create(&permission.PermissionActionRel{
			PermissionID:permissionID,
			ActionID:v,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &actionIDs, nil
}
