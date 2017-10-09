package soda

import "maizuo.com/soda/erp/api/src/server/model"

type Account struct {
	model.Model
	UserID   int	`json:"userId,omitempty"`
	App      string `json:"app,omitempty"`
	Key      string `json:"key,omitempty"`
	Extra    string `json:"extra,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Password string	`json:"password"`
	Status   int	`json:"status"`
}

func (Account) TableName() string {
	return "account"
}
