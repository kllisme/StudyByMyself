package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type City struct {
	model.Model
	Name         string        `json:"name"`
	Code         string        `json:"code"`
	ParentCode   string        `json:"parentCode"`
	ProvinceName string        `json:"provinceName" gorm:"-"`
	ProvinceCode string        `json:"provinceCode" gorm:"-"`
}

func (City) TableName() string {
	return "city"
}
