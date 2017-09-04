package public

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"time"
	"github.com/spf13/viper"
	"github.com/kataras/iris/core/errors"
)

type Advertisement struct {
	model.Model
	Name            string        `json:"name"`
	Title           string        `json:"title"`
	URL             string        `json:"url"`
	LocationID      int        `json:"locationId"`
	LocationName    string        `json:"locationName" gorm:"-"`
	APPID           int        `json:"appId" gorm:"-"`
	APPName         string        `json:"appName" gorm:"-"`
	StartedAt       time.Time        `json:"startedAt"`
	EndedAt         time.Time        `json:"endedAt"`
	Status          int        `json:"status"`
	Image           string        `json:"image"`
	Order           int        `json:"order"`
	DisplayStrategy int        `json:"displayStrategy"`
	DisplayParams   string        `json:"displayParams"`
}

func (Advertisement) TableName() string {
	return "ad"
}

func (a *Advertisement) AfterFind() error  {
	if a.Image != "" {
		fullURL := viper.GetString("resource.oss.domain")+"/"+viper.GetString("resource.oss.object.ad")+a.Image
		a.Image = fullURL
		return nil
	}
	return errors.New("图片链接为空")
}
