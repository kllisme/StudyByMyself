package soda

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type Ticket struct {
	model.Model
	TicketId          string `json:"ticketId"`
	UserId            int    `json:"userId"`
	Mobile            string `json:"mobile"`
	BillId            string `json:"billId"`
	Value             int    `json:"value"`
	Token             string `json:"token"`
	SnapShot          string `json:"snapShot"`
	DeviceSerial      string `json:"deviceSerial"`
	DeviceMode        int    `json:"deviceMode"`
	OwnerId           int    `json:"ownerId"`
	APPID             int        `json:"appId"`
	Status            int    `json:"status"`
	DeviceReferenceID int        `json:"deviceReferenceId"`
	PaymentId         int    `json:"paymentId"`
	Feature           int        `json:"feature"`
	CreatedTimestamp  int        `json:"createdTimestamp"`
	Settle            bool   `json:"settle",gorm:"-"`
}

func (Ticket) TableName() string {
	return "ticket"
}

