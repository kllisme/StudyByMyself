package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type School struct {
	model.Model
	Name         string        `json:"name"`
	ProvinceCode string        `json:"provinceCode"`
	ProvinceName string        `json:"provinceName" gorm:"-"`
	CityCode     string        `json:"cityCode"`
	CityName     string        `json:"cityName"  gorm:"-"`
	AreaCode     string        `json:"areaCode"`
	AreaName     string        `json:"areaName"  gorm:"-"`
}

func (School) TableName() string {
	return "school"
}
