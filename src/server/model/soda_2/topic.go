package soda_2

import (
	"time"
	"github.com/spf13/viper"
	"maizuo.com/soda/erp/api/src/server/entity"
	"encoding/json"
	"maizuo.com/soda/erp/api/src/server/model"
)

type Topic struct {
	model.Model
	Title            string        `json:"title"`              //名称
	Description      string        `json:"description"`
	Content          string        `json:"content"`
	Value            int        `json:"value"`
	UserID           int        `json:"userId"`
	UserName         string        `json:"userName" gorm:"-"`
	CityID           int        `json:"cityId"`
	CityName         string        `json:"cityName"`
	SchoolID         int        `json:"schoolId"`
	SchoolName       string        `json:"schoolName"`
	ChannelID        int        `json:"channelId"`
	ChannelTitle     string        `json:"channelTitle"`
	StartedAt        *time.Time        `json:"startedAt"`
	EndedAt          *time.Time        `json:"endedAt"`
	Images           string        `json:"images"`
	HasOnline        int        `json:"hasOnline"`
	UniqueVisitor    int        `json:"uniqueVisitor"`
	Likes            int        `json:"likes"`
	Consultation     int        `json:"consultation" gorm:"-"` //商品询问人数
	Comments         int        `json:"comments"`
	Status           int        `json:"status"`
	CreatedTimestamp int        `json:"createdTimestamp"`
}

func (Topic) TableName() string {
	return "topic"
}

func (t *Topic) AfterFind() error {
	if t.Images != "" {
		imageList := make([]*entity.Image, 0)
		if err := json.Unmarshal([]byte(t.Images), &imageList); err != nil {
			return err
		}

		for _, image := range imageList {
			image.URL = viper.GetString("resource.oss.domain") + "/" + viper.GetString("resource.oss.object.topic") + image.URL
		}
		buff, err := json.Marshal(imageList)
		t.Images = string(buff)
		if err != nil {
			return err
		}
	}
	return nil
}
