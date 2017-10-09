package soda_2

import "maizuo.com/soda/erp/api/src/server/model"

type Channel struct {
	model.Model
	Title            string `json:"title"` //名称
	URL              string        `json:"url"`
	Order            int    `json:"order"`
	IsOfficial       int        `json:"isOfficial"`
	Status           int        `json:"status"`
	CreatedTimestamp int        `json:"createdTimestamp"`
}

func (Channel) TableName() string {
	return "channel"
}
