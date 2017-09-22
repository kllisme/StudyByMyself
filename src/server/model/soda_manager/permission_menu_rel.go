package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type PermissionMenuRel struct {
	model.Model
	PermissionID int `json:"permissionId"`
	MenuID       int        `json:"menuId"`
}

func (PermissionMenuRel)TableName() string {
	return "erp_permission_menu_rel"
}
