package soda_2

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	twoModel "maizuo.com/soda/erp/api/src/server/model/soda_2"
)

type UserService struct {

}

func (self *UserService)GetByID(id int) (*twoModel.User, error) {
	user := twoModel.User{}
	err := common.Soda2DB_R.Where(id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (self *UserService)Count(cityID int) (int, error) {
	count := 0
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if cityID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("city_id = ?", cityID)
		})
	}
	err := common.Soda2DB_R.Table("user").Scopes(scopes...).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (self *UserService)FilterUserIDs(name string) ([]int, error) {
	userIDs := make([]int,0)
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like (?)", "%" + name + "%")
		})
	}
	err := common.Soda2DB_R.Table("user").Scopes(scopes...).Pluck("id",&userIDs).Error
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}
