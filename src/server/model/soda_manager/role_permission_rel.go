package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type RolePermissionRel struct {
	model.Model
	RoleID       int        `json:"roleId"`
	PermissionID int        `json:"permissionId"`
}

func (RolePermissionRel) TableName() string {
	return "erp_role_permission_rel"
}
