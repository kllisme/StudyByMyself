package soda

import "maizuo.com/soda/erp/api/src/server/model"

type Payment struct {
	model.Model
	Name        string        `json:"name"`
	Description string        `json:"description"`
	LogoURL     string        `json:"logoUrl"`
	Icon        string        `json:"icon"`
	Tag         string        `json:"tag"`
	Status      int        `json:"status"`
}

func (Payment) TableName() string {
	return "payment"
}
