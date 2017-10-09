package soda

import "maizuo.com/soda/erp/api/src/server/model"

type CBRel struct {
	model.Model
	ChipcardId     int        `json:"chipcardId,omitempty"`
	BusinessUserId int        `json:"businessUserId,omitempty"`
	UserId         int        `json:"userId,omitempty"`
	status         int        `json:"status"`
}

func (CBRel) TableName() string {
	return "chipcard_business_rel"
}
