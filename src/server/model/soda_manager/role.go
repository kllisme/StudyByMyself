package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type Role struct {
	model.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

func (Role) TableName() string {
	return "erp_role"
}
