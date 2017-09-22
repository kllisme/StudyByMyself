package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type BillRel struct {
	model.Model
	BillId      string `json:"billId"`
	BatchNo     string `json:"batchNo"`
	Type        int    `json:"type"`
	IsSuccessed bool   `json:"isSuccessed"`
	Reason      string `json:"reason"`
	OuterNo     string `json:"outerNo"`
	BillType    int    `json:"billType"`
	ErrCode     string `json:"errCode"`
}

func (BillRel) TableName() string {
	return "bill_rel"
}
