package permission

import "maizuo.com/soda/erp/api/src/server/model"

type UserRoleRel struct {
	model.Model
	UserID 	int	`json:"userId"`
	RoleID	int 	`json:"roleId"`
}

func (UserRoleRel)TableName() string {
	return "erp_user_role_rel"
}
