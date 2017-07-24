package permission

import "maizuo.com/soda/erp/api/src/server/model"

type Action struct {
	model.Model
	HandlerName string        `json:"handlerName"` //处理函数名
	Method      string        `json:"method"`
	API         string        `json:"api"`
	Description string        `json:"description"`
}

func (Action) TableName() string {
	return "erp_action"
}
