package permission

import (
	"maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/entity"
)

type PermissionService struct {

}

func (self *PermissionService)GetByID(id int) (*permission.Permission, error) {
	permission := permission.Permission{}
	err := common.SodaMngDB_R.Where(id).Find(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (self *PermissionService)GetAll() (*[]*permission.Permission, error) {
	permissionList := make([]*permission.Permission, 0)
	if err := common.SodaMngDB_R.Order("id desc").Find(&permissionList).Error; err != nil {
		return nil, err
	}
	return &permissionList, nil

}

func (self *PermissionService)Paging(categoryID int, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	permissionList := make([]*permission.Permission, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if categoryID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("category_id = ?", categoryID)
		})
	}
	if err := db.Model(&permission.Permission{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&permissionList).Error; err != nil {
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

func (self *PermissionService)GetListByIDs(ids ...interface{}) (*[]*permission.Permission, error) {
	permissionList := make([]*permission.Permission, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Order("id desc").Find(&permissionList).Error
	if err != nil {
		return nil, err
	}
	return &permissionList, nil
}

func (self *PermissionService)Create(permission *permission.Permission) (*permission.Permission, error) {
	err := common.SodaMngDB_WR.Create(&permission).Error
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (self *PermissionService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&permission.Permission{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("permission_id = ?", id).Delete(&permission.RolePermissionRel{}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("permission_id = ?", id).Delete(&permission.PermissionActionRel{}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("permission_id = ?", id).Delete(&permission.PermissionElementRel{}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("permission_id = ?", id).Delete(&permission.PermissionMenuRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *PermissionService)Update(entity *permission.Permission) (*permission.Permission, error) {
	_permission := map[string]interface{}{
		"name":entity.Name,
		"category_id":entity.CategoryID,
	}
	if err := common.SodaMngDB_WR.Model(&permission.Permission{}).Where(entity.ID).Updates(_permission).Scan(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
