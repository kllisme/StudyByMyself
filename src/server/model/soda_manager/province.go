package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type Province struct {
	model.Model
	Name string `json:"name"`
	Code string `json:"code"`
}

func (Province) TableName() string {
	return "province"
}
