package soda_manager

import (
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
	"github.com/jinzhu/gorm"
)

type BillRelService struct {
}

func (self *BillRelService) Create(billRelList ...*mngModel.BillRel) (int, error) {
	var err error
	var r *gorm.DB
	var isTrue bool
	rows := 0
	tx := common.SodaMngDB_WR.Begin()
	for _, billRel := range billRelList {
		isTrue = tx.NewRecord(&billRel)
		if !isTrue {
			e := &functions.DefinedError{}
			e.Msg = "can not create a new record!"
			err = e
			tx.Rollback()
			return 0, err
		}
		r = tx.Create(&billRel)
		if r.Error != nil {
			tx.Rollback()
			return 0, r.Error
		}else {
			rows += int(r.RowsAffected)
		}
	}
	tx.Commit()
	return rows, nil
}
