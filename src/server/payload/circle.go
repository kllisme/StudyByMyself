package payload

type Circle struct {
	Order       int        `json:"order"`
	CityID      int        `json:"cityId"`
	CityName    string        `json:"cityName"`
	SchoolCount int        `json:"schoolCount"`
	UserCount   int        `json:"userCount"`
	TopicCount  int        `json:"topicCount"`
}

