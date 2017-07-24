package permission

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type PermissionActionRel struct {
	model.Model
	PermissionID int        `json:"permissionId"`
	ActionID     int        `json:"actionId"`
}

func (PermissionActionRel)TableName() string {
	return "erp_permission_action_rel"
}
