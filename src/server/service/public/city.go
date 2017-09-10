package public

import (
	"maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type CityService struct {

}

func (self *CityService)GetByID(id int) (*public.City, error) {
	city := public.City{}
	err := common.SodaMngDB_R.Where(id).Find(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (self *CityService)GetByCode(code string) (*public.City, error) {
	city := public.City{}
	err := common.SodaMngDB_R.Where("code = ?", code).First(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (self *CityService)Paging(name string, parentCode string, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	cityList := make([]*public.City, 0)
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
	if err := db.Model(&public.City{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&cityList).Error; err != nil {
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

func (self *CityService)Create(city *public.City) (*public.City, error) {
	err := common.SodaMngDB_WR.Create(&city).Error
	if err != nil {
		return nil, err
	}
	return city, nil
}

func (self *CityService)Update(city *public.City) (*public.City, error) {
	_city := map[string]interface{}{
		"name":city.Name,
		"code":city.Code,
		"parent_code":city.ParentCode,
	}
	if err := common.SodaMngDB_WR.Model(&public.City{}).Where(city.ID).Updates(_city).Scan(city).Error; err != nil {
		return nil, err
	}
	return city, nil
}

func (self *CityService)Delete(id int) error {
	city := public.City{}
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Model(&public.City{}).Where(id).Scan(&city).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Delete(&public.City{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Model(&public.Area{}).Where("parent_code = ?", city.Code).Update("parent_code", "").Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
