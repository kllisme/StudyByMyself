package permission

import "maizuo.com/soda/erp/api/src/server/model"

type Menu struct {
	model.Model
	Name     string	`json:"name"`
	Icon     string `json:"icon"`
	Url      string	`json:"url"`
	ParentID int    `json:"parentId"`
	Level    int    `json:"level"`
	Status   int    `json:"status"`
	Position int    `json:"position"`
}

func (Menu) TableName() string {
	return "erp_menu"
}
