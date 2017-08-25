package community

type Channel struct {
	Model
	Title            string `json:"title"` //名称
	URL              string        `json:"url"`
	Order            int    `json:"order"`
	IsOfficial       int        `json:"is_official"`
	Status           int        `json:"status"`
	CreatedTimestamp int        `json:"created_timestamp"`
}

func (Channel) TableName() string {
	return "2_channel"
}
