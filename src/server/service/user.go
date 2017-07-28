package service

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/entity"
)

type UserService struct {
}

//校验所有信息
func (self *UserService) CheckInfo(user *model.User) (*model.User, error) {
	db := common.SodaMngDB_R
	result := model.User{}
	if err := db.Where(user).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *UserService)Paging(name string, account string, id int, roleID int, page int, perPage int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	userList := make([]*model.User, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like (?)", "%" + name + "%")
		})
	}
	if account != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("account like (?)", "%" + account + "%")
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
	if err := db.Model(&model.User{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset((page - 1) * perPage).Limit(perPage).Find(&userList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = (page - 1) * perPage
	pagination.Pagination.To = perPage * page
	if pagination.Pagination.To > pagination.Pagination.Total {
		pagination.Pagination.To = pagination.Pagination.Total
	}
	pagination.Objects = userList
	return &pagination, nil

}

func (self *UserService) GetById(id int) (*model.User, error) {
	db := common.SodaMngDB_R
	user := model.User{}
	if err := db.Where(id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (self *UserService) GetByAccount(account string) (*model.User, error) {
	db := common.SodaMngDB_R
	result := model.User{}
	if err := db.Where(&model.User{Account:account}).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *UserService) GetList(limit int, offset int) (interface{}, error) {

	return nil, nil
}

func (self *UserService) Create(user *model.User) (*model.User, error) {
	err := common.SodaMngDB_WR.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (self *UserService) UpdateById(user *model.User) (*model.User, error) {
	if err := common.SodaMngDB_R.Model(&model.User{}).Updates(&user).Where(user.ID).Scan(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (self *UserService) UpdateByMobile(in interface{}) (interface{}, error) {
	return nil, nil
}

func (self *UserService) DeleteById(id int) (interface{}, error) {
	return nil, nil
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
