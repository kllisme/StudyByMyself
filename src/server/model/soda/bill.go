package soda

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"time"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type Bill struct {
	model.Model
	BillID           string	`json:"billId"`
	UserID           int	`json:"userId"`
	Mobile           string `json:"mobile"`
	WalletID         int	`json:"billId"`
	Value            int    `json:"value"`
	Title            string `json:"title"`
	Type             int	`json:"type"`
	Tradeno          string	`json:"tradeno"`
	Action           int	`json:"action"`
	CreatedTimestamp int 	`json:"createdTimestamp"`
	OldID            int    `json:"oldId"`
	Status           int    `json:"status"`
}

func (Bill) TableName() string {
	return "bill"
}

func (self *Bill) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now().Local()
	at := now.Format("2006-01-02 15:04:05")
	scope.SetColumn("created_at", at)
	scope.SetColumn("updated_at", at)
	scope.SetColumn("created_timestamp", now.Unix())
	scope.SetColumn("bill_id", functions.GenerateIdByMobile(self.Mobile))
	return nil
}
