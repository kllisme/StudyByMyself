package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type Area struct {
	model.Model
	Name         string        `json:"name"`
	Code         string        `json:"code"`
	ParentCode   string        `json:"parentCode"`
	ProvinceName string        `json:"provinceName" gorm:"-"`
	ProvinceCode string        `json:"provinceCode" gorm:"-"`
	CityName     string        `json:"cityName" gorm:"-"`
	CityCode     string        `json:"cityCode" gorm:"-"`
}

func (Area) TableName() string {
	return "area"
}
