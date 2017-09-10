package two

import (
	"time"
)

type User struct {
	Model
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
