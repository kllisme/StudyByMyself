package permission

import "maizuo.com/soda/erp/api/src/server/model"

type API struct {
	model.Model
	Name       string        `json:"name"`  //处理函数名
	Status     int        `json:"status"`
	Title      string        `json:"title"` //描述性标题
	CategoryID int        `json:"categoryId"`
}

func (API) TableName() string {
	return "erp_api"
}
