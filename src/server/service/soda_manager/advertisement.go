package soda_manager

import (
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type AdvertisementService struct {

}

func (self *AdvertisementService)GetListByADPositionID(id int) (*[]*mngModel.Advertisement, error) {
	advertisementList := make([]*mngModel.Advertisement, 0)
	err := common.SodaMngDB_R.Where("ad_position_id = ?", id).Order("id desc").Find(&advertisementList).Error
	if err != nil {
		return nil, err
	}
	return &advertisementList, nil
}

func (self *AdvertisementService)GetByID(id int) (*mngModel.Advertisement, error) {
	advertisement := mngModel.Advertisement{}
	err := common.SodaMngDB_R.Where(id).Find(&advertisement).Error
	if err != nil {
		return nil, err
	}
	return &advertisement, nil
}

func (self *AdvertisementService)Paging(fullName string, name string, adPositionIDs []int, start string, end string, display int, status int, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	advertisementList := make([]*mngModel.Advertisement, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if fullName != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name = ?", fullName)
		})
	}

	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like (?)", "%" + name + "%")
		})
	}
	if len(adPositionIDs) != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("ad_position_id in (?)", adPositionIDs)
		})
	}
	if start != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("started_at >= ?", start)
		})
	}
	if end != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("ended_at <= ?", end)
		})
	}
	if display != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("display_strategy = ?", display)
		})
	}
	if status != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", status)
		})
	}

	if err := db.Model(&mngModel.Advertisement{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&advertisementList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = advertisementList
	return &pagination, nil

}

func (self *AdvertisementService)Create(advertisement *mngModel.Advertisement) (*mngModel.Advertisement, error) {
	err := common.SodaMngDB_WR.Create(&advertisement).Error
	if err != nil {
		return nil, err
	}
	return advertisement, nil
}

func (self *AdvertisementService)Update(advertisement *mngModel.Advertisement) (*mngModel.Advertisement, error) {
	_advertisement := map[string]interface{}{
		"name":advertisement.Name,
		"adPositionId":advertisement.AdPositionID,
		"title":advertisement.Title,
		"image":advertisement.Image,
		"url":advertisement.URL,
		"startedAt":advertisement.StartedAt,
		"endedAt":advertisement.EndedAt,
		"displayStrategy":advertisement.DisplayStrategy,
		"displayParams":advertisement.DisplayParams,
		"status":advertisement.Status,
	}
	if err := common.SodaMngDB_WR.Model(&mngModel.Advertisement{}).Where(advertisement.ID).Updates(_advertisement).Scan(advertisement).Error; err != nil {
		return nil, err
	}
	return advertisement, nil
}

func (self *AdvertisementService)BatchUpdateOrder(advertisementList *[]*mngModel.Advertisement) (*[]*mngModel.Advertisement, error) {
	result := []*mngModel.Advertisement{}
	for _, ad := range *advertisementList {
		advertisement := &mngModel.Advertisement{}
		if err := common.SodaMngDB_WR.Model(&mngModel.Advertisement{}).Where(ad.ID).Update("order", ad.Order).Scan(advertisement).Error; err != nil {
			return nil, err
		}
		result = append(result, advertisement)
	}
	return &result, nil
}

func (self *AdvertisementService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Delete(&mngModel.Advertisement{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

