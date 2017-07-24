package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type RolePermissionRelService struct {

}

func (self *RolePermissionRelService) GetPermissionIDsByRoleIDs(roleIDs ...interface{}) ([]int, error) {
	permissionIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.RolePermissionRel{}).Where("role_id in (?)", roleIDs...).Pluck("permission_id",&permissionIDs).Error
	if err != nil {
		return nil, err
	}
	return permissionIDs, nil
}
