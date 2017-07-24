package permission

import (
	"maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
)

type ActionService struct {

}

func (self *ActionService)GetListByIds(ids ...interface{}) (*[]*permission.Action, error) {
	actionList := make([]*permission.Action, 0)
	err := common.SodaMngDB_R.Where("id in (?)",ids...).Find(&actionList).Error
	if err != nil {
		return nil, err
	}
	return &actionList, nil
}

