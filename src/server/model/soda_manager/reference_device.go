package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type ReferenceDevice struct {
	model.Model
	Name string	`json:"name"`
}

func (ReferenceDevice) TableName() string {
	return "reference_device"
}
