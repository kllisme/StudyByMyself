package soda

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"time"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type ChipcardBill struct {
	model.Model
	BillId           string        `json:"billId,omitempty"`
	Value            int        `json:"value,omitempty"`
	Mobile           string        `json:"mobile"`
	UserId           int        `json:"userId,omitempty"`
	ChipcardId       int        `json:"chipcardId,omitempty"`
	Title            string        `json:"title,omitempty"`
	Type             int        `json:"type"`
	Action           int        `json:"action"`
	Status           int        `json:"status"`
	CreatedTimestamp int        `json:"created_timestamp,omitempty"`
}

func (self *ChipcardBill) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now().Local()
	at := now.Format("2006-01-02 15:04:05")
	scope.SetColumn("created_at", at)
	scope.SetColumn("updated_at", at)
	scope.SetColumn("created_timestamp", now.Unix())
	scope.SetColumn("bill_id", functions.GenerateIdByMobile(self.Mobile))
	return nil
}

func (ChipcardBill) TableName() string {
	return "chipcard_bill"
}
