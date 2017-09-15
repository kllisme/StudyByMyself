package soda

import (
	"maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/common"
)

type PaymentService struct {
}


func (self *PaymentService) GetAll() (*[]*soda.Payment, error) {
	paymentList := []*soda.Payment{}
	err := common.SodaDB_R.Find(&paymentList).Error
	if err != nil {
		return nil, err
	}
	return &paymentList, nil
}
