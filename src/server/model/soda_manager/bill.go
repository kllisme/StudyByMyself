package soda_manager

import (
	"time"
	"maizuo.com/soda/erp/api/src/server/model"
)

type Bill struct {
	model.Model
	BillId      string    `json:"billId"`
	Account     string    `json:"account"`
	AccountType int       `json:"accountType"`
	AccountName string    `json:"accountName"`
	UserId      int       `json:"userId"`
	UserName    string    `json:"userName"`
	UserAccount string    `json:"userAccount"`
	SettledAt   time.Time `json:"settledAt"`
	TotalAmount int       `json:"totalAmount"`
	Count       int       `json:"count"`
	Rate        int       `json:"rate"`
	Cast        int       `json:"cast"`
	Amount      int       `json:"amount"`
	SubmittedAt time.Time `json:"submittedAt"`
	Mobile      string    `json:"mobile"`
	RealName    string    `json:"realName"`
	BankName    string    `json:"bankName"`
	Status      int       `json:"status"`
	Mode        int       `json:"mode"`
}

func (Bill) TableName() string {
	return "bill"
}
