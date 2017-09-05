package public

import (
	"maizuo.com/soda/erp/api/src/server/model"
)

type ADSpace struct {
	model.Model
	Name           string        `json:"name"`
	APPID          int        `json:"appId"`
	APPName        string       `gorm:"-" json:"appName"`
	IdentifyNeeded int        `json:"identifyNeeded"`
	Description    string        `json:"description"`
	Standard       string        `json:"standard"`
}

func (ADSpace) TableName() string {
	return "ad_space"
}

