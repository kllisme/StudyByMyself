package permission

import (
	"maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
)

type APIService struct {

}

func (self *APIService)GetListByIds(ids ...interface{}) (*[]permission.API, error) {
	apiList := make([]permission.API, 0)
	err := common.SodaMngDB_R.Where("id in (?)",ids...).Find(&apiList).Error
	if err != nil {
		return nil, err
	}
	return &apiList, nil
}

