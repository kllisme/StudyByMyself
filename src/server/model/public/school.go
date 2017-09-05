package public

import "maizuo.com/soda/erp/api/src/server/model"

type School struct {
	model.Model
	Name        string `json:"name"`
	ProvinceId  int    `json:"provinceId"`
}

func (School) TableName() string {
	return "school"
}
