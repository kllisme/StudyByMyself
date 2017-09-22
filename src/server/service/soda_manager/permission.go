package soda_manager

import (
	mngService "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/entity"
)

type PermissionService struct {

}

func (self *PermissionService)GetByID(id int) (*mngService.Permission, error) {
	permission := mngService.Permission{}
	err := common.SodaMngDB_R.Where(id).Find(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (self *PermissionService)GetAll() (*[]*mngService.Permission, error) {
	permissionList := make([]*mngService.Permission, 0)
	if err := common.SodaMngDB_R.Order("id desc").Find(&permissionList).Error; err != nil {
		return nil, err
	}
	return &permissionList, nil

}

func (self *PermissionService)Paging(categoryID int, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	permissionList := make([]*mngService.Permission, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if categoryID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("category_id = ?", categoryID)
		})
	}
	if err := db.Model(&mngService.Permission{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&permissionList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = permissionList
	return &pagination, nil
}

func (self *PermissionService)GetListByIDs(ids ...interface{}) (*[]*mngService.Permission, error) {
	permissionList := make([]*mngService.Permission, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Order("id desc").Find(&permissionList).Error
	if err != nil {
		return nil, err
	}
	return &permissionList, nil
}

func (self *PermissionService)Create(permission *mngService.Permission) (*mngService.Permission, error) {
	err := common.SodaMngDB_WR.Create(&permission).Error
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (self *PermissionService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&mngService.Permission{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("permission_id = ?", id).Delete(&mngService.RolePermissionRel{}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("permission_id = ?", id).Delete(&mngService.PermissionActionRel{}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("permission_id = ?", id).Delete(&mngService.PermissionElementRel{}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("permission_id = ?", id).Delete(&mngService.PermissionMenuRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *PermissionService)Update(entity *mngService.Permission) (*mngService.Permission, error) {
	_permission := map[string]interface{}{
		"name":entity.Name,
		"category_id":entity.CategoryID,
	}
	if err := common.SodaMngDB_WR.Model(&mngService.Permission{}).Where(entity.ID).Updates(_permission).Scan(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
