package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type RoleMenuRelService struct {

}

func (self *RoleMenuRelService) GetMenuIDsByRoleIDs(roleIDs ...interface{}) ([]int, error) {
	menuIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.RoleMenuRel{}).Where("role_id in (?)", roleIDs...).Pluck("menu_id",&menuIDs).Error
	if err != nil {
		return nil, err
	}
	return menuIDs, nil
}
