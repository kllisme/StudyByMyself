package permission

import "maizuo.com/soda/erp/api/src/server/model"

type Permission struct {
	model.Model
	Name       string        `json:"name"`
	//Type       int        `json:"type"`
	Status     int        `json:"status"`
	CategoryID int        `json:"categoryId"`
}

func (Permission)TableName() string {
	return "erp_permission"
}
