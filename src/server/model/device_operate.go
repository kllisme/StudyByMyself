package model

type DeviceOperate struct {
	Model
	OperatorID   int    `json:"operatorId"`
	OperatorType int    `json:"operatorType"`
	SerialNumber string `json:"serialNumber"`
	UserID       int    `json:"userId"`
	FromUserID   int    `json:"fromUserId"`
	ToUserID     int    `json:"toUserId"`
}

func (DeviceOperate) TableName() string {
	return "device_operate"
}
