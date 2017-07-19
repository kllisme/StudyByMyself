package permission

import "maizuo.com/soda/erp/api/src/server/model"

type RoleOperation struct {
	model.Model
	Operator int        `json:"operator"`
	RoleID   int        `json:"roleId"`
}

func (RoleOperation)TableName() string {
	return "erp_user_role_rel"
}
