package permission

import "maizuo.com/soda/erp/api/src/server/model"

type Menu struct {
	model.Model
	Name     string `json:"name"`
	Icon     string	`json:"icon"`
	Url      string `json:"url"`
	ParentID int    `json:"pid"`
	Level    int    `json:"level"`
	Status   int    `json:"status"`
}

func (Menu) TableName() string {
	return "erp_menu"
}
