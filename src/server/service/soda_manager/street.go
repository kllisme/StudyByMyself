package soda_manager

import (
	mngService "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type StreetService struct {

}

func (self *StreetService)GetByID(id int) (*mngService.Street, error) {
	street := mngService.Street{}
	err := common.SodaMngDB_R.Where(id).Find(&street).Error
	if err != nil {
		return nil, err
	}
	return &street, nil
}

func (self *StreetService)GetByCode(code string) (*mngService.Street, error) {
	street := mngService.Street{}
	err := common.SodaMngDB_R.Where("code = ?", code).First(&street).Error
	if err != nil {
		return nil, err
	}
	return &street, nil
}

func (self *StreetService)Paging(name string, parentCode string, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	streetList := make([]*mngService.Street, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like (?)", "%" + name + "%")
		})
	}
	if parentCode != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_code = ?", parentCode)
		})
	}
	if err := db.Model(&mngService.Street{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&streetList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = streetList
	return &pagination, nil

}

func (self *StreetService)Create(street *mngService.Street) (*mngService.Street, error) {
	err := common.SodaMngDB_WR.Create(&street).Error
	if err != nil {
		return nil, err
	}
	return street, nil
}

func (self *StreetService)Update(street *mngService.Street) (*mngService.Street, error) {
	_street := map[string]interface{}{
		"name":street.Name,
		"code":street.Code,
		"parent_code":street.ParentCode,
	}
	if err := common.SodaMngDB_WR.Model(&mngService.Street{}).Where(street.ID).Updates(_street).Scan(street).Error; err != nil {
		return nil, err
	}
	return street, nil
}

func (self *StreetService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&mngService.Street{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
