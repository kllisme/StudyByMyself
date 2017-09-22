package soda_manager

import (
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type MenuService struct {

}

func (self *MenuService)GetListByIDs(ids ...interface{}) (*[]*mngModel.Menu, error) {
	menuList := make([]*mngModel.Menu, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Order("position").Find(&menuList).Error
	if err != nil {
		return nil, err
	}
	return &menuList, nil
}

func (self *MenuService)GetByID(id int) (*mngModel.Menu, error) {
	menu := mngModel.Menu{}
	err := common.SodaMngDB_R.Where(id).Find(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (self *MenuService)Paging(offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	menuList := make([]*mngModel.Menu, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if err := db.Model(&mngModel.Menu{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("position").Find(&menuList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = menuList
	return &pagination, nil

}

func (self *MenuService)Create(menu *mngModel.Menu) (*mngModel.Menu, error) {
	err := common.SodaMngDB_WR.Create(&menu).Error
	if err != nil {
		return nil, err
	}
	return menu, nil
}

func (self *MenuService)Update(menu *mngModel.Menu) (*mngModel.Menu, error) {
	if err := common.SodaMngDB_WR.Model(&mngModel.Menu{}).Save(&menu).Where(menu.ID).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (self *MenuService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&mngModel.Menu{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("menu_id = ?", id).Delete(&mngModel.PermissionMenuRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *MenuService)RearrangePosition(menuList *[]*mngModel.Menu) (*[]*mngModel.Menu, error) {
	result := []*mngModel.Menu{}
	for _, menu := range *menuList {
		//menu := &permission.Menu{}
		if err := common.SodaMngDB_WR.Model(&mngModel.Menu{}).Where(menu.ID).Update("position", menu.Position).Scan(menu).Error; err != nil {
			return nil, err
		}
		result = append(result, menu)
	}
	return &result, nil
}
