package community

import "time"

type Inbox struct {
	Model
	Title            string	`json:"title"` //名称
	Content          string        `json:"content"`
	Icon             string        `json:"icon"`
	IsOfficial       int        `json:"is_official"`
	Type             int        `json:"type"`
	OperatorID       int        `json:"operator_id"`
	Action           int        `json:"action"`
	TargetType       int        `json:"target_type"`
	TargetID         int        `json:"target_id"`
	SenderID         int        `json:"sender_id"`
	ReceiverID       int        `json:"receiver_id"`
	SenderUserID     int        `json:"sender_user_id"`
	ReceiverUserID   int        `json:"receiver_user_id"`
	Extra            string        `json:"extra"`
	Status           int        `json:"status"`
	CreatedTimestamp int        `json:"created_timestamp"`
	UpdatedTimestamp *time.Time        `json:"updated_timestamp"`
}

func (Inbox) TableName() string {
	return "inbox"
}
