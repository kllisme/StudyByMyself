package soda

import (
	"maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type ChipcardBillService  struct {

}

func (self *ChipcardBillService) Paging(mobile string, action int, _type int, start string, end string, limit int, offset int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	chipcardBillList := make([]*soda.ChipcardBill, 0)
	db := common.SodaDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if mobile != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("mobile = ?", mobile)
		})
	}

	if action != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("action = ?", action)
		})
	}

	if _type != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", _type)
		})
	}

	if start != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("created_timestamp >= UNIX_TIMESTAMP(?)", start)
		})
	}
	if end != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("created_timestamp < UNIX_TIMESTAMP(DATE_ADD(?, INTERVAL 1 day))", end)
		})
	}

	if err := db.Model(&soda.ChipcardBill{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("created_timestamp desc").Find(&chipcardBillList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = chipcardBillList
	return &pagination, nil
}
