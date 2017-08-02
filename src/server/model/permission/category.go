package permission

import "maizuo.com/soda/erp/api/src/server/model"

type Category struct {
	model.Model
	Name 	string 	`json:"name"`
	Status 	int	`json:"status"`
}

func (Category) TableName() string {
	return "erp_category"
}
