package public

import (
	"maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type AdvertisementService struct {

}

func (self *AdvertisementService)GetListByLocationID(id int) (*[]*public.Advertisement, error) {
	advertisementList := make([]*public.Advertisement, 0)
	err := common.SodaMngDB_R.Where("location_id = ?", id).Order("id desc").Find(&advertisementList).Error
	if err != nil {
		return nil, err
	}
	return &advertisementList, nil
}

func (self *AdvertisementService)GetByID(id int) (*public.Advertisement, error) {
	advertisement := public.Advertisement{}
	err := common.SodaMngDB_R.Where(id).Find(&advertisement).Error
	if err != nil {
		return nil, err
	}
	return &advertisement, nil
}

func (self *AdvertisementService)Paging(title string,locationIDs []int, start string,end string,display int,status int, page int, perPage int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	advertisementList := make([]*public.Advertisement, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if title != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("title like (?)", "%" + title + "%")
		})
	}
	if len(locationIDs) != 0 {
		scopes = append(scopes,func(db *gorm.DB) *gorm.DB {
			return db.Where("location_id in (?)", locationIDs)
		})
	}
	if start != "" {
		scopes = append(scopes,func(db *gorm.DB) *gorm.DB {
			return db.Where("started_at >= ?", start)
		})
	}
	if end != "" {
		scopes = append(scopes,func(db *gorm.DB) *gorm.DB {
			return db.Where("ended_at <= ?", end)
		})
	}
	if display != 0 {
		scopes = append(scopes,func(db *gorm.DB) *gorm.DB {
			return db.Where("display_strategy = ?", display)
		})
	}
	if status != 0 {
		scopes = append(scopes,func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", status)
		})
	}


	if err := db.Model(&public.Advertisement{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset((page - 1) * perPage).Limit(perPage).Order("id desc").Find(&advertisementList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = (page - 1) * perPage + 1
	pagination.Pagination.To = perPage * page
	if pagination.Pagination.To > pagination.Pagination.Total {
		pagination.Pagination.To = pagination.Pagination.Total
	}
	pagination.Objects = advertisementList
	return &pagination, nil

}

func (self *AdvertisementService)Create(advertisement *public.Advertisement) (*public.Advertisement, error) {
	err := common.SodaMngDB_WR.Create(&advertisement).Error
	if err != nil {
		return nil, err
	}
	return advertisement, nil
}

func (self *AdvertisementService)Update(advertisement *public.Advertisement) (*public.Advertisement, error) {
	_advertisement := map[string]interface{}{
		"name":advertisement.Name,
		"locationId":advertisement.LocationID,
		"title":advertisement.Title,
		"image":advertisement.Image,
		"url":advertisement.URL,
		"startedAt":advertisement.StartedAt,
		"endedAt":advertisement.EndedAt,
		"displayStrategy":advertisement.DisplayStrategy,
		"displayParams":advertisement.DisplayParams,
		"status":advertisement.Status,
	}
	if err := common.SodaMngDB_WR.Model(&public.Advertisement{}).Where(advertisement.ID).Updates(_advertisement).Scan(advertisement).Error; err != nil {
		return nil, err
	}
	return advertisement, nil
}

func (self *AdvertisementService)BatchUpdateOrder(advertisementList *[]*public.Advertisement) (*[]*public.Advertisement, error) {
	result := []*public.Advertisement{}
	for _, ad := range *advertisementList {
		advertisement := &public.Advertisement{}
		if err := common.SodaMngDB_WR.Model(&public.Advertisement{}).Where(ad.ID).Update("order", ad.Order).Scan(advertisement).Error; err != nil {
			return nil, err
		}
		result = append(result, advertisement)
	}
	return &result, nil
}

func (self *AdvertisementService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Delete(&public.Advertisement{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

