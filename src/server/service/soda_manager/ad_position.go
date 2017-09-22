package soda_manager

import (
 	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type ADPositionService struct {

}

func (self *ADPositionService)GetByID(id int) (*mngModel.ADPosition, error) {
	adPosition := mngModel.ADPosition{}
	err := common.SodaMngDB_R.Where(id).Find(&adPosition).Error
	if err != nil {
		return nil, err
	}
	return &adPosition, nil
}

func (self *ADPositionService)Paging(name string, appID int, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	adPositionList := make([]*mngModel.ADPosition, 0)
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

	if err := db.Model(&mngModel.ADPosition{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&adPositionList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = adPositionList
	return &pagination, nil

}

func (self *ADPositionService)Create(adPosition *mngModel.ADPosition) (*mngModel.ADPosition, error) {
	err := common.SodaMngDB_WR.Create(&adPosition).Error
	if err != nil {
		return nil, err
	}
	return adPosition, nil
}

func (self *ADPositionService)Update(adPosition *mngModel.ADPosition) (*mngModel.ADPosition, error) {
	_adPosition := map[string]interface{}{
		"name":adPosition.Name,
		"description":adPosition.Description,
		"identifyNeeded": adPosition.IdentifyNeeded,
		"appId":adPosition.APPID,
		"standard":adPosition.Standard,
	}
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Model(&mngModel.ADPosition{}).Where(adPosition.ID).Updates(_adPosition).Scan(adPosition).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if adPosition.IdentifyNeeded == 0 {
		_ad := map[string]interface{}{
			"display_strategy":1,
			"display_params":"",
		}
		if err := tx.Model(&mngModel.Advertisement{}).Where("ad_position_id = ?", adPosition.ID).Updates(_ad).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return adPosition, nil
}

func (self *ADPositionService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Delete(&mngModel.ADPosition{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Where("ad_position_id = ?", id).Delete(&mngModel.Advertisement{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *ADPositionService)GetIDsByAPPID(appID int) ([]int, error) {
	adPositionIDs := make([]int, 0)
	if err := common.SodaMngDB_R.Model(&mngModel.ADPosition{}).Where("app_id = ?", appID).Pluck("id", &adPositionIDs).Error; err != nil {
		return nil, err
	}

	return adPositionIDs, nil
}
