package soda

import (
	"maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/common"
)

type WalletService  struct {

}

func (self *WalletService) GetByUserMobile(mobile string) (*soda.Wallet, error) {
	db := common.SodaDB_R
	result := soda.Wallet{}
	if err := db.Where(&soda.Wallet{Mobile:mobile}).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
