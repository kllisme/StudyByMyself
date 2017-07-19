package permission

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type RoleAPIRel struct {
	model.Model
	RoleID int        `json:"roleId"`
	APIID  int        `json:"apiId"`
}

func (RoleAPIRel)TableName() string {
	return "erp_role_api_rel"
}
