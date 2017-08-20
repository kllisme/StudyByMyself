package permission

import (
	"maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type MenuService struct {

}

func (self *MenuService)GetListByIDs(ids ...interface{}) (*[]*permission.Menu, error) {
	menuList := make([]*permission.Menu, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Order("id desc").Find(&menuList).Error
	if err != nil {
		return nil, err
	}
	return &menuList, nil
}

func (self *MenuService)GetByID(id int) (*permission.Menu, error) {
	menu := permission.Menu{}
	err := common.SodaMngDB_R.Where(id).Find(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (self *MenuService)Paging(page int, perPage int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	menuList := make([]*permission.Menu, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if err := db.Model(&permission.Menu{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset((page - 1) * perPage).Limit(perPage).Order("id desc").Find(&menuList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = (page - 1) * perPage
	pagination.Pagination.To = perPage * page
	if pagination.Pagination.To > pagination.Pagination.Total {
		pagination.Pagination.To = pagination.Pagination.Total
	}
	pagination.Objects = menuList
	return &pagination, nil

}

func (self *MenuService)Create(menu *permission.Menu) (*permission.Menu, error) {
	err := common.SodaMngDB_WR.Create(&menu).Error
	if err != nil {
		return nil, err
	}
	return menu, nil
}

func (self *MenuService)Update(menu *permission.Menu) (*permission.Menu, error) {
	if err := common.SodaMngDB_WR.Model(&permission.Menu{}).Save(&menu).Where(menu.ID).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (self *MenuService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&permission.Menu{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("menu_id = ?", id).Delete(&permission.PermissionMenuRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
