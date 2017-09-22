package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type PermissionElementRel struct {
	model.Model
	PermissionID int `json:"permissionId"`
	ElementID       int        `json:"elementId"`
}

func (PermissionElementRel)TableName() string {
	return "erp_permission_element_rel"
}
