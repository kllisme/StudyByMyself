package permission

import "maizuo.com/soda/erp/api/src/server/model"

type Element struct {
	model.Model
	Name      string        `json:"name"`
	Reference string        `json:"reference"`
}

func (Element)TableName() string {
	return "erp_element"
}
