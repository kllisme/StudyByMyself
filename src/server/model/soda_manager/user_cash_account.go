package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type UserCashAccount struct {
	model.Model
	UserId       int    `json:"userId"`
	Type         int    `json:"type"`
	RealName     string `json:"realName"`
	HeadBankName string `json:"headBankName"`
	BankName     string `json:"bankName"`
	Account      string `json:"account"`
	Mobile       string `json:"mobile"`
	CityId       int    `json:"cityId"`
	ProvinceId   int    `json:"provinceId"`
	Mode         int    `json:"mode"`
}

func (UserCashAccount) TableName() string {
	return "user_cash_account"
}
