package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type GroupRoleRel struct {
	model.Model
	GroupID 	int	`json:"groupId"`
	RoleID	int 	`json:"roleId"`
}

func (GroupRoleRel)TableName() string {
	return "erp_group_role_rel"
}
