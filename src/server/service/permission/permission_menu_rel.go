package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type PermissionMenuRelService struct {

}

func (self *PermissionMenuRelService) GetMenuIDsByPermissionIDs(permissionIDs ...interface{}) ([]int, error) {
	menuIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.PermissionMenuRel{}).Where("permission_id in (?)", permissionIDs...).Pluck("menu_id", &menuIDs).Error
	if err != nil {
		return nil, err
	}
	return menuIDs, nil
}
