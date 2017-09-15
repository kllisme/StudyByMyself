package soda

import "maizuo.com/soda/erp/api/src/server/model"

type Bill struct {
	model.Model
	BillID           string `json:"billId"`
	UserID           int    `json:"userId"`
	Mobile           string `json:"mobile"`
	WalletID         int `json:"billId"`
	Value            int    `json:"value"`
	Title            string `json:"title"`
	Type             int        `json:"type"`
	Tradeno          string        `json:"tradeno"`
	Action           int        `json:"action"`
	CreatedTimestamp int        `json:"createdTimestamp"`
	OldID            int    `json:"oldId"`
	Status           int    `json:"status"`
}

func (Bill) TableName() string {
	return "bill"
}

