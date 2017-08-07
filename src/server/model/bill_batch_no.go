package model

type BillBatchNo struct {
	Model
	BillId   string `json:"billId"`
	BatchNo  string `json:"batchNo"`
	BillType int    `json:"billType"`
}

func (BillBatchNo) TableName() string {
	return "bill_batch_no"
}
