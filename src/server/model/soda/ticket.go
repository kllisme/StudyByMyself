package soda

import (
	"time"
)

type Ticket struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	TicketId     string `json:"ticketId"`
	UserId       int    `json:"userId"`
	Mobile       string `json:"mobile"`
	BillId       string `json:"billId"`
	Value        int    `json:"value"`
	Token        string `json:"token"`
	SnapShot     string `json:"snapShot"`
	DeviceSerial string `json:"deviceSerial"`
	DeviceMode   int    `json:"deviceMode"`
	OwnerId      int    `json:"ownerId"`
	Status       int    `json:"status"`
	PaymentId    int    `json:"paymentId"`
	Settle       bool   `json:"settle",gorm:"-"`
}

func (Ticket) TableName() string {
	return "ticket"
}

