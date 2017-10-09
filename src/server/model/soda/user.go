package soda

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"time"
)

type User struct {
	model.Model
	NickName         string        `json:"nickName"`
	AvatorURL        string                `json:"avatorUrl"`
	Mobile           string        `json:"mobile"`
	Gender           int                `json:"gender"`
	Status           int         `json:"status"`
	CreatedTimestamp int          `json:"createdTimestamp"`
	DeletedAt        *time.Time        `json:"deletedAt"`
}

func (User) TableName() string {
	return "user"
}
