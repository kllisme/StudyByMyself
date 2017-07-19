package service

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"maizuo.com/soda/erp/api/src/server/common"
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

func (self *UserService) GetById(id int) (*model.User, error) {
	db := common.SodaMngDB_R
	result := model.User{}
	if err := db.Where(id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
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
	err:=common.SodaMngDB_WR.Create(user).Error
	if err != nil {
		return nil,err
	}
	return user, nil
}

func (self *UserService) UpdateById(in interface{}) (interface{}, error) {
	return nil, nil
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
