package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type PermissionActionRelService struct {

}

func (self *PermissionActionRelService) GetActionIDsByPermissionIDs(permissionIDs ...interface{}) ([]int, error) {
	actionIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.PermissionActionRel{}).Where("permission_id in (?)", permissionIDs...).Pluck("action_id",&actionIDs).Error
	if err != nil {
		return nil, err
	}
	return actionIDs, nil
}
