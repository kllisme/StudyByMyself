package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type Application struct {
	model.Model
	Name        string	`json:"name"`
	Description string	`json:"description"`
}

func (Application) TableName() string {
	return "application"
}
