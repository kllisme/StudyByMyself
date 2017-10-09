package soda_2

import (
	"time"
	"maizuo.com/soda/erp/api/src/server/model"
)

type Inbox struct {
	model.Model
	Title            string	`json:"title"` //名称
	Content          string        `json:"content"`
	Icon             string        `json:"icon"`
	IsOfficial       int        `json:"isOfficial"`
	Type             int        `json:"type"`
	OperatorID       int        `json:"operatorId"`
	Action           int        `json:"action"`
	TargetType       int        `json:"targetType"`
	TargetID         int        `json:"targetId"`
	SenderID         int        `json:"senderId"`
	ReceiverID       int        `json:"receiverId"`
	SenderUserID     int        `json:"senderUserId"`
	ReceiverUserID   int        `json:"receiverUserId"`
	Extra            string        `json:"extra"`
	Status           int        `json:"status"`
	CreatedTimestamp int        `json:"createdTimestamp"`
	UpdatedTimestamp *time.Time        `json:"updatedTimestamp"`
}

func (Inbox) TableName() string {
	return "inbox"
}
