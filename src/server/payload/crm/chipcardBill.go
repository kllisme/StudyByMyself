package crm

import "time"

type ChipCardBill struct {
	ActionName string        `json:"actionName"`
	Value      string        `json:"value"`
	Time       time.Time	 `json:"time"`
	Operator   string        `json:"operator"`
}
