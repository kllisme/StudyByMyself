package permission

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/permission"
)

type UserRoleRelService struct {

}

func (self *UserRoleRelService) GetRoleIDsByUserID(userID int) ([]int, error) {
	roleIDs := make([]int, 0)
	err := common.SodaMngDB_R.Model(&permission.UserRoleRel{}).Where("user_id = ?", userID).Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}

func (self *UserRoleRelService) AsignRoles(userID int, roleIDs []int) (interface{}, error) {
	tx := common.SodaMngDB_R.Begin()
	err := tx.Delete(permission.UserRoleRel{}, "user_id = ?", userID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range roleIDs {
		err := tx.Create(permission.UserRoleRel{
			UserID:userID,
			RoleID:v,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return nil, nil
}
