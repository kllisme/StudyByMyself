package soda_2

import (
	"time"
	"maizuo.com/soda/erp/api/src/server/model"
)

type User struct {
	model.Model
	Name             string `json:"name"`   //名称
	OpenID           string        `json:"openId"`
	AvatorURL        string        `json:"avatorUrl"`
	Gender           int    `json:"gender"` //性别
	Bio              string `json:"bio"`
	SchoolID         int        `json:"schoolId"`
	Birthday         *time.Time `json:"birthday"`
	Status           int        `json:"status"`
	CreatedTimestamp int        `json:"createdTimestamp"`
}

func (User) TableName() string {
	return "user"
}
