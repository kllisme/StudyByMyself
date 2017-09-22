package soda_manager

import (
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
)

type RoleService struct {

}

func (self *RoleService)GetByID(id int) (*mngModel.Role, error) {
	role := mngModel.Role{}
	err := common.SodaMngDB_R.Where(id).Find(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (self *RoleService)GetAll() (*[]*mngModel.Role, error) {
	roleList := make([]*mngModel.Role, 0)
	if err := common.SodaMngDB_R.Order("id desc").Find(&roleList).Error; err != nil {
		return nil, err
	}
	return &roleList, nil

}

func (self *RoleService)GetListByIDs(ids ...interface{}) (*[]*mngModel.Role, error) {
	roleList := make([]*mngModel.Role, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Order("id desc").Find(&roleList).Error
	if err != nil {
		return nil, err
	}
	return &roleList, nil
}

func (self *RoleService)Create(role *mngModel.Role) (*mngModel.Role, error) {
	err := common.SodaMngDB_WR.Create(&role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (self *RoleService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&mngModel.Role{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("role_id = ?", id).Delete(&mngModel.UserRoleRel{}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("role_id = ?", id).Delete(&mngModel.RolePermissionRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *RoleService)Update(role *mngModel.Role) (*mngModel.Role, error) {
	_role := map[string]interface{}{
		"name":role.Name,
		"description":role.Description,
	}
	if err := common.SodaMngDB_WR.Model(&mngModel.Role{}).Where(role.ID).Updates(_role).Scan(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}
