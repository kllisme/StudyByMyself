package soda

import "maizuo.com/soda/erp/api/src/server/model"

type Chipcard struct {
	model.Model
	Value      int        `json:"value,omitempty"`
	Mobile     string     `json:"mobile,omitempty"`
	UserId     int        `json:"userId,omitempty"`
	Status     int        `json:"status"`
	OperatorId int        `json:"operator_id,omitempty"`
}

func (Chipcard) TableName() string {
	return "chipcard"
}
