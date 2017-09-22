package soda_manager

import (
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type UserService struct {
}

//校验所有信息
func (self *UserService) CheckInfo(user *mngModel.User) (*mngModel.User, error) {
	db := common.SodaMngDB_R
	result := mngModel.User{}
	if err := db.Where(user).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *UserService) Paging(name string, account string, id int, roleID int, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	userList := make([]*mngModel.User, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like (?)", "%"+name+"%")
		})
	}
	if account != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("account like (?)", "%"+account+"%")
		})
	}
	if id != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where(id)
		})
	}
	if roleID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Joins("left join erp_user_role_rel rel on rel.user_id = user.id and rel.role_id = ?", roleID)
		})
	}
	if err := db.Model(&mngModel.User{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&userList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = userList
	return &pagination, nil

}

func (self *UserService) GetById(id int) (*mngModel.User, error) {
	db := common.SodaMngDB_R
	user := mngModel.User{}
	if err := db.Where(id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (self *UserService) GetByAccount(account string) (*mngModel.User, error) {
	db := common.SodaMngDB_R
	result := mngModel.User{}
	if err := db.Unscoped().Where(&mngModel.User{Account:account}).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *UserService) GetList(limit int, offset int) (interface{}, error) {

	return nil, nil
}

func (self *UserService) Create(user *mngModel.User) (*mngModel.User, error) {
	err := common.SodaMngDB_WR.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (self *UserService) Update(user *mngModel.User) (*mngModel.User, error) {
	_user := map[string]interface{}{
		"name":      user.Name,
		"contact":   user.Contact,
		"mobile":    user.Mobile,
		"telephone": user.Telephone,
		"address":   user.Address,
		"email":     user.Email,
	}
	if err := common.SodaMngDB_WR.Model(&mngModel.User{}).Where("id = ?", user.ID).Updates(_user).Scan(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (self *UserService) UpdateByMobile(in interface{}) (interface{}, error) {
	return nil, nil
}

func (self *UserService) DeleteById(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Delete(&mngModel.User{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("user_id = ?", id).Delete(&mngModel.UserRoleRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *UserService) ChangePassword(id int, password string) (*mngModel.User, error) {
	user := mngModel.User{}
	err := common.SodaMngDB_WR.Model(&mngModel.User{}).Where(id).Updates(&mngModel.User{Password: password}).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (self *UserService) DeleteByMobile(mobile string) (interface{}, error) {
	return nil, nil
}

func (self *UserService) RemoveById(id int) (interface{}, error) {
	return nil, nil
}

func (self *UserService) RemoveByMobile(mobile string) (interface{}, error) {
	return nil, nil
}
