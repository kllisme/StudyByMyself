package model

type BillRel struct {
	Model
	BillId      string `json:"billId"`
	BatchNo     string `json:"batchNo"`
	Type        int    `json:"type"`
	IsSuccessed bool   `json:"isSuccessed"`
	Reason      string `json:"reason"`
	OuterNo     string `json:"outerNo"`
	BillType    int    `json:"billType"`
}

func (BillRel) TableName() string {
	return "bill_rel"
}
