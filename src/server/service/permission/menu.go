package permission

import (
	"maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
)

type MenuService struct {

}

func (self *MenuService)GetListByIds(ids ...interface{}) (*[]*permission.Menu, error) {
	menuList := make([]*permission.Menu, 0)
	err := common.SodaMngDB_R.Where("id in (?)",ids...).Find(&menuList).Error
	if err != nil {
		return nil, err
	}
	return &menuList, nil
}
