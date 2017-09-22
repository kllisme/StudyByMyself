package soda_manager

import "maizuo.com/soda/erp/api/src/server/model"

type BillBatchNo struct {
	model.Model
	BillId   string `json:"billId"`
	BatchNo  string `json:"batchNo"`
	BillType int    `json:"billType"`
}

func (BillBatchNo) TableName() string {
	return "bill_batch_no"
}
