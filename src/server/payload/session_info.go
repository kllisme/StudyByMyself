package payload

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type SessionInfo struct {
	User     *model.User        `json:"user"`
	MenuList *[]permission.Menu       `json:"menuList"`
	APIList	*[]permission.API	`json:"apiList"`
}
