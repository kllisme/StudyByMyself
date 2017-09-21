package public

import (
	"maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type StreetService struct {

}

func (self *StreetService)GetByID(id int) (*public.Street, error) {
	street := public.Street{}
	err := common.SodaMngDB_R.Where(id).Find(&street).Error
	if err != nil {
		return nil, err
	}
	return &street, nil
}

func (self *StreetService)GetByCode(code string) (*public.Street, error) {
	street := public.Street{}
	err := common.SodaMngDB_R.Where("code = ?", code).First(&street).Error
	if err != nil {
		return nil, err
	}
	return &street, nil
}

func (self *StreetService)Paging(name string, parentCode string, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	streetList := make([]*public.Street, 0)
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
	if err := db.Model(&public.Street{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&streetList).Error; err != nil {
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

func (self *StreetService)Create(street *public.Street) (*public.Street, error) {
	err := common.SodaMngDB_WR.Create(&street).Error
	if err != nil {
		return nil, err
	}
	return street, nil
}

func (self *StreetService)Update(street *public.Street) (*public.Street, error) {
	_street := map[string]interface{}{
		"name":street.Name,
		"code":street.Code,
		"parent_code":street.ParentCode,
	}
	if err := common.SodaMngDB_WR.Model(&public.Street{}).Where(street.ID).Updates(_street).Scan(street).Error; err != nil {
		return nil, err
	}
	return street, nil
}

func (self *StreetService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&public.Street{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
