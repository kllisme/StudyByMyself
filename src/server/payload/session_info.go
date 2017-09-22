package payload

import (
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type SessionInfo struct {
	User        *mngModel.User        `json:"user"`
	MenuList    *[]*mngModel.Menu       `json:"menuList"`
	ElementList *[]*mngModel.Element        `json:"elementList"`
}
