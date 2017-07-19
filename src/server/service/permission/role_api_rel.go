package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type RoleAPIRelService struct {

}

func (self *RoleAPIRelService) GetAPIIDsByRoleIDs(roleIDs ...interface{}) ([]int, error) {
	apiIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.RoleAPIRel{}).Where("role_id in (?)", roleIDs...).Pluck("api_id",&apiIDs).Error
	if err != nil {
		return nil, err
	}
	return apiIDs, nil
}
