package soda_manager

import (
	manModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type ProvinceService struct {

}

func (self *ProvinceService)GetByID(id int) (*manModel.Province, error) {
	province := manModel.Province{}
	err := common.SodaMngDB_R.Where(id).Find(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func (self *ProvinceService)GetByCode(code string) (*manModel.Province, error) {
	province := manModel.Province{}
	err := common.SodaMngDB_R.Where("code = ?", code).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func (self *ProvinceService)Paging(offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	provinceList := make([]*manModel.Province, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if err := db.Model(&manModel.Province{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&provinceList).Error; err != nil {
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

func (self *ProvinceService)Create(province *manModel.Province) (*manModel.Province, error) {
	err := common.SodaMngDB_WR.Create(&province).Error
	if err != nil {
		return nil, err
	}
	return province, nil
}

func (self *ProvinceService)Update(province *manModel.Province) (*manModel.Province, error) {
	_province := map[string]interface{}{
		"name":province.Name,
		"code":province.Code,
	}
	if err := common.SodaMngDB_WR.Model(&manModel.Province{}).Where(province.ID).Updates(_province).Scan(province).Error; err != nil {
		return nil, err
	}
	return province, nil
}

func (self *ProvinceService)Delete(id int) error {
	province := manModel.Province{}
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Model(&manModel.Province{}).Where(id).Scan(&province).Error; err!=nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Delete(&manModel.Province{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Model(&manModel.City{}).Where("parent_code = ?", province.Code).Update("parent_code", "").Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Model(&manModel.School{}).Where("province_code = ?", province.Code).Update("province_code", "").Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
