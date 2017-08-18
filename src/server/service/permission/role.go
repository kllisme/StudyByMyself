package permission

import (
	"maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
)

type RoleService struct {

}

func (self *RoleService)GetByID(id int) (*permission.Role, error) {
	role := permission.Role{}
	err := common.SodaMngDB_R.Where(id).Find(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (self *RoleService)GetAll() (*[]*permission.Role, error) {
	roleList := make([]*permission.Role, 0)
	if err := common.SodaMngDB_R.Find(&roleList).Error; err != nil {
		return nil, err
	}
	return &roleList, nil

}

func (self *RoleService)GetListByIDs(ids ...interface{}) (*[]*permission.Role, error) {
	roleList := make([]*permission.Role, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Find(&roleList).Error
	if err != nil {
		return nil, err
	}
	return &roleList, nil
}

func (self *RoleService)Create(role *permission.Role) (*permission.Role, error) {
	err := common.SodaMngDB_WR.Create(&role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (self *RoleService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&permission.Role{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("role_id = ?", id).Delete(&permission.UserRoleRel{}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("role_id = ?", id).Delete(&permission.RolePermissionRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *RoleService)Update(role *permission.Role) (*permission.Role, error) {
	_role := map[string]interface{}{
		"name":role.Name,
		"description":role.Description,
	}
	if err := common.SodaMngDB_WR.Model(&permission.Role{}).Where(role.ID).Updates(_role).Scan(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}
