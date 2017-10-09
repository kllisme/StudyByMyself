package soda

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"github.com/jinzhu/gorm"
	"time"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type ChipcardRecharge struct {
	model.Model
	RechargeId       string        `json:"recharge_id,omitempty"`
	UserId           int        `json:"user_id,omitempty"`
	Mobile           string        `json:"mobile,omitempty"`
	OperatorId       int        `json:"operator_id,omitempty"`
	ChipcardId       int        `json:"chipcard_id,omitempty"`
	BillId           string        `json:"bill_id,omitempty"`
	Snapshot         string        `json:"snapshot,omitempty"`
	Count            int        `json:"count,omitempty"` //默认为1
	Value            int        `json:"value,omitempty"`
	Extra            string        `json:"extra"`
	Status           int        `json:"status,omitempty"` //7:已发货
	CreatedTimestamp int        `json:"created_timestamp,omitempty"`
}

func (self *ChipcardRecharge) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now().Local()
	at := now.Format("2006-01-02 15:04:05")
	scope.SetColumn("created_at", at)
	scope.SetColumn("updated_at", at)
	scope.SetColumn("created_timestamp", now.Unix())
	scope.SetColumn("recharge_id", functions.GenerateIdByMobile(self.Mobile))
	return nil
}

func (ChipcardRecharge) TableName() string {
	return "chipcard_recharge"
}
