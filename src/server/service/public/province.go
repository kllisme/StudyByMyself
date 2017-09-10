package public

import (
	"maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type ProvinceService struct {

}

func (self *ProvinceService)GetByID(id int) (*public.Province, error) {
	province := public.Province{}
	err := common.SodaMngDB_R.Where(id).Find(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func (self *ProvinceService)GetByCode(code string) (*public.Province, error) {
	province := public.Province{}
	err := common.SodaMngDB_R.Where("code = ?", code).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func (self *ProvinceService)Paging(offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	provinceList := make([]*public.Province, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if err := db.Model(&public.Province{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&provinceList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = provinceList
	return &pagination, nil

}

func (self *ProvinceService)Create(province *public.Province) (*public.Province, error) {
	err := common.SodaMngDB_WR.Create(&province).Error
	if err != nil {
		return nil, err
	}
	return province, nil
}

func (self *ProvinceService)Update(province *public.Province) (*public.Province, error) {
	_province := map[string]interface{}{
		"name":province.Name,
		"code":province.Code,
	}
	if err := common.SodaMngDB_WR.Model(&public.Province{}).Where(province.ID).Updates(_province).Scan(province).Error; err != nil {
		return nil, err
	}
	return province, nil
}

func (self *ProvinceService)Delete(id int) error {
	province := public.Province{}
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Model(&public.Province{}).Where(id).Scan(&province).Error; err!=nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Delete(&public.Province{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Model(&public.City{}).Where("parent_code = ?", province.Code).Update("parent_code", "").Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
