package public

import (
	"maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type AreaService struct {

}

func (self *AreaService)GetByID(id int) (*public.Area, error) {
	area := public.Area{}
	err := common.SodaMngDB_R.Where(id).Find(&area).Error
	if err != nil {
		return nil, err
	}
	return &area, nil
}

func (self *AreaService)GetByCode(code string) (*public.Area, error) {
	area := public.Area{}
	err := common.SodaMngDB_R.Where("code = ?", code).First(&area).Error
	if err != nil {
		return nil, err
	}
	return &area, nil
}

func (self *AreaService)Paging(name string, parentCode string, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	areaList := make([]*public.Area, 0)
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
	if err := db.Model(&public.Area{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&areaList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = areaList
	return &pagination, nil

}

func (self *AreaService)Create(area *public.Area) (*public.Area, error) {
	err := common.SodaMngDB_WR.Create(&area).Error
	if err != nil {
		return nil, err
	}
	return area, nil
}

func (self *AreaService)Update(area *public.Area) (*public.Area, error) {
	_area := map[string]interface{}{
		"name":area.Name,
		"code":area.Code,
		"parent_code":area.ParentCode,
	}
	if err := common.SodaMngDB_WR.Model(&public.Area{}).Where(area.ID).Updates(_area).Scan(area).Error; err != nil {
		return nil, err
	}
	return area, nil
}

func (self *AreaService)Delete(id int) error {
	area := public.Area{}
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Model(&public.Area{}).Where(id).Scan(&area).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Delete(&public.Area{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Model(&public.Street{}).Where("parent_code = ?", area.Code).Update("parent_code", "").Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
