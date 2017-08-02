package permission

import "maizuo.com/soda/erp/api/src/server/model"

type Group struct {
	model.Model
	Name       string        `json:"name"`
}

func (Group)TableName() string {
	return "erp_group"
}
