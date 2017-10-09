package soda

import (
	"maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/common"
)

type ChipcardService  struct {

}

func (self *ChipcardService) GetByUserMobile(mobile string) (*soda.Chipcard, error) {
	db := common.SodaDB_R
	result := soda.Chipcard{}
	if err := db.Where(&soda.Chipcard{Mobile:mobile}).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
