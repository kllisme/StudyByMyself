package public

import "maizuo.com/soda/erp/api/src/server/model"

type Street struct {
	model.Model
	Name         string        `json:"name"`
	Code         string        `json:"code"`
	ParentCode   string        `json:"parentCode"`
	ProvinceName string        `json:"provinceName" gorm:"-"`
	ProvinceCode string        `json:"provinceCode" gorm:"-"`
	CityName     string        `json:"cityName" gorm:"-"`
	CityCode     string        `json:"cityCode" gorm:"-"`
	AreaName     string        `json:"areaName" gorm:"-"`
	AreaCode     string        `json:"areaCode" gorm:"-"`
}

func (Street) TableName() string {
	return "street"
}
