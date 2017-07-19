package permission

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type RoleMenuRel struct {
	model.Model
	RoleID	int `json:"roleId"`
	MenuID	int	`json:"menuId"`
}

func (RoleMenuRel)TableName() string {
	return "erp_role_menu_rel"
}
