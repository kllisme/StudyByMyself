package public

import (
	"maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type ADSpaceService struct {

}

func (self *ADSpaceService)GetByID(id int) (*public.ADSpace, error) {
	adSpace := public.ADSpace{}
	err := common.SodaMngDB_R.Where(id).Find(&adSpace).Error
	if err != nil {
		return nil, err
	}
	return &adSpace, nil
}

func (self *ADSpaceService)Paging(name string, appID int, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	adSpaceList := make([]*public.ADSpace, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if appID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("app_id = ?", appID)
		})
	}

	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name = ?", name)
		})
	}

	if err := db.Model(&public.ADSpace{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&adSpaceList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = adSpaceList
	return &pagination, nil

}

func (self *ADSpaceService)Create(adSpace *public.ADSpace) (*public.ADSpace, error) {
	err := common.SodaMngDB_WR.Create(&adSpace).Error
	if err != nil {
		return nil, err
	}
	return adSpace, nil
}

func (self *ADSpaceService)Update(adSpace *public.ADSpace) (*public.ADSpace, error) {
	_adSpace := map[string]interface{}{
		"name":adSpace.Name,
		"description":adSpace.Description,
		"identifyNeeded": adSpace.IdentifyNeeded,
		"appId":adSpace.APPID,
		"standard":adSpace.Standard,
	}
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Model(&public.ADSpace{}).Where(adSpace.ID).Updates(_adSpace).Scan(adSpace).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if adSpace.IdentifyNeeded == 0 {
		_ad := map[string]interface{}{
			"display_strategy":1,
			"display_params":"",
		}
		if err := tx.Model(&public.Advertisement{}).Where("location_id = ?", adSpace.ID).Updates(_ad).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return adSpace, nil
}

func (self *ADSpaceService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Delete(&public.ADSpace{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Where("location_id = ?", id).Delete(&public.Advertisement{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *ADSpaceService)GetLocationIDs(appID int) ([]int, error) {
	locationIDs := make([]int, 0)
	if err := common.SodaMngDB_R.Model(&public.ADSpace{}).Where("app_id = ?", appID).Pluck("id", &locationIDs).Error; err != nil {
		return nil, err
	}

	return locationIDs, nil
}
