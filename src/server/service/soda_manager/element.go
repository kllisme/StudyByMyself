package soda_manager

import (
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type ElementService struct {

}

func (self *ElementService)GetListByIDs(ids ...interface{}) (*[]*mngModel.Element, error) {
	elementList := make([]*mngModel.Element, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Order("id desc").Find(&elementList).Error
	if err != nil {
		return nil, err
	}
	return &elementList, nil
}

func (self *ElementService)GetByID(id int) (*mngModel.Element, error) {
	element := mngModel.Element{}
	err := common.SodaMngDB_R.Where(id).Find(&element).Error
	if err != nil {
		return nil, err
	}
	return &element, nil
}

func (self *ElementService)Paging(name string, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	elementList := make([]*mngModel.Element, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like (?)", "%" + name + "%")
		})
	}
	if err := db.Model(&mngModel.Element{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&elementList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = elementList
	return &pagination, nil

}

func (self *ElementService)Create(element *mngModel.Element) (*mngModel.Element, error) {
	err := common.SodaMngDB_WR.Create(&element).Error
	if err != nil {
		return nil, err
	}
	return element, nil
}

func (self *ElementService)Update(element *mngModel.Element) (*mngModel.Element, error) {
	_element := map[string]interface{}{
		"name":element.Name,
		"reference":element.Reference,
	}
	if err := common.SodaMngDB_WR.Model(&mngModel.Element{}).Where(element.ID).Updates(_element).Scan(element).Error; err != nil {
		return nil, err
	}
	return element, nil
}

func (self *ElementService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&mngModel.Element{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("element_id = ?", id).Delete(&mngModel.PermissionElementRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
