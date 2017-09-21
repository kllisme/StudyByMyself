package public

import (
	"maizuo.com/soda/erp/api/src/server/common"
	model "maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type SchoolService struct {

}

func (self *SchoolService)GetByID(id int) (*model.School, error) {
	school := model.School{}
	err := common.SodaMngDB_R.Where(id).Find(&school).Error
	if err != nil {
		return nil, err
	}
	return &school, nil
}

func (self *SchoolService)Paging(name string, provinceID string, cityID string, areaID string, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	schoolList := make([]*model.School, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like (?)", "%" + name + "%")
		})
	}
	if provinceID != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("province_code = ?", provinceID)
		})
	}
	if cityID != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("city_code = ?", cityID)
		})
	}
	if areaID != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("area_code = ?", areaID)
		})
	}

	if err := db.Model(&model.School{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&schoolList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = schoolList
	return &pagination, nil

}

func (self *SchoolService)Create(school *model.School) (*model.School, error) {
	err := common.SodaMngDB_WR.Create(&school).Error
	if err != nil {
		return nil, err
	}
	return school, nil
}

func (self *SchoolService)Update(school *model.School) (*model.School, error) {
	_school := map[string]interface{}{
		"name":school.Name,
		"province_code":school.ProvinceCode,
		"city_code":school.CityCode,
		"area_code":school.AreaCode,
	}
	if err := common.SodaMngDB_WR.Model(&model.School{}).Where(school.ID).Updates(_school).Scan(school).Error; err != nil {
		return nil, err
	}
	return school, nil
}

func (self *SchoolService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&model.School{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

