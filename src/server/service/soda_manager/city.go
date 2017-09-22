package soda_manager

import (
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type CityService struct {

}

func (self *CityService)GetByID(id int) (*mngModel.City, error) {
	city := mngModel.City{}
	err := common.SodaMngDB_R.Where(id).Find(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (self *CityService)GetByCode(code string) (*mngModel.City, error) {
	city := mngModel.City{}
	err := common.SodaMngDB_R.Where("code = ?", code).First(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (self *CityService)Paging(name string, parentCode string, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	cityList := make([]*mngModel.City, 0)
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
	if err := db.Model(&mngModel.City{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&cityList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = cityList
	return &pagination, nil

}

func (self *CityService)Create(city *mngModel.City) (*mngModel.City, error) {
	err := common.SodaMngDB_WR.Create(&city).Error
	if err != nil {
		return nil, err
	}
	return city, nil
}

func (self *CityService)Update(city *mngModel.City) (*mngModel.City, error) {
	_city := map[string]interface{}{
		"name":city.Name,
		"code":city.Code,
		"parent_code":city.ParentCode,
	}
	if err := common.SodaMngDB_WR.Model(&mngModel.City{}).Where(city.ID).Updates(_city).Scan(city).Error; err != nil {
		return nil, err
	}
	return city, nil
}

func (self *CityService)Delete(id int) error {
	city := mngModel.City{}
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Model(&mngModel.City{}).Where(id).Scan(&city).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Delete(&mngModel.City{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Model(&mngModel.Area{}).Where("parent_code = ?", city.Code).Update("parent_code", "").Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Model(&mngModel.School{}).Where("city_code = ?", city.Code).Update("city_code", "").Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
