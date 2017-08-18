package model

import "time"

type DailyBill struct {
	Model
	UserId      int       `json:"userId"`
	UserName    string    `json:"userName"`
	TotalAmount int       `json:"totalAmount"`
	SettledAt   time.Time `json:"settledAt"`
	BillAt      string    `json:"billAt"`
	Status      int       `json:"status"` //0:未申请提现or未结账 1:已申请提现 2:已结账 3:结账中 4:结账失败
	SubmitAt    time.Time `json:"submitAt"`
	OrderCount  int       `json:"orderCount"`
	AccountType int       `json:"accountType"`
	AccountName string    `json:"accountName"`
	Account     string    `json:"account"`
	RealName    string    `json:"realName"`
	BankName    string    `json:"bankName"`
	Mobile      string    `json:"mobile"`
	HasMarked   int       `json:"hasMarked"`
	BillId      string    `json:"billId"`
}

func (DailyBill) TableName() string {
	return "daily_bill"
}
