package community

import (
	"time"
)

type Topic struct {
	Model
	Title            string        `json:"title"` //名称
	Description      string        `json:"description"`
	Content          string        `json:"content"`
	Value            int        `json:"value"`
	UserID           int        `json:"userId"`
	UserName         string        `json:"userName"`
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
	Comments         int        `json:"comments"`
	Status           int        `json:"status"`
	CreatedTimestamp int        `json:"createdTimestamp"`
}

func (Topic) TableName() string {
	return "2_topic"
}
