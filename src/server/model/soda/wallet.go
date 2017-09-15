package soda

import "maizuo.com/soda/erp/api/src/server/model"

type Wallet struct {
	model.Model
	Value  int        `json:"value"`
	UserID int    `json:"userId"`
	Mobile string `json:"mobile"`
	Status int    `json:"status"`
}

func (Wallet) TableName() string {
	return "wallet"
}

