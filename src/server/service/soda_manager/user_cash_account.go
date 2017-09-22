package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/common"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type UserCashAccountService struct {
}

func (self *UserCashAccountService) Basic(id int) (*mngModel.UserCashAccount, error) {
	userCashAccount := &mngModel.UserCashAccount{}
	r := common.SodaMngDB_R.Where("id = ?", id).First(userCashAccount)
	if r.Error != nil {
		return nil, r.Error
	}
	return userCashAccount, nil
}

func (self *UserCashAccountService) BasicByUserId(userId int) (*mngModel.UserCashAccount, error) {
	userCashAccount := &mngModel.UserCashAccount{}
	r := common.SodaMngDB_R.Where("user_id = ?", userId).First(userCashAccount)
	if r.Error != nil {
		return nil, r.Error
	}
	return userCashAccount, nil
}

func (self *UserCashAccountService) BasicMapByUserId(userIds interface{}) (map[int]*mngModel.UserCashAccount, error) {
	list := &[]*mngModel.UserCashAccount{}
	accountMap := make(map[int]*mngModel.UserCashAccount)
	r := common.SodaMngDB_R.Where("user_id in (?)", userIds).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, account := range *list {
		accountMap[account.UserId] = account
	}
	return accountMap, nil
}

func (self *UserCashAccountService) BasicMapByType(payType ...int) (map[int]*mngModel.UserCashAccount, error) {
	list := &[]*mngModel.UserCashAccount{}
	accountMap := make(map[int]*mngModel.UserCashAccount, 0)
	r := common.SodaMngDB_R.Where("type in (?)", payType).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, account := range *list {
		accountMap[account.UserId] = account
	}
	return accountMap, nil
}

func (self *UserCashAccountService) Create(userCashAccount *mngModel.UserCashAccount) (bool,error) {
	transAction := common.SodaMngDB_WR.Begin()
	r := transAction.Create(userCashAccount)
	if r.RowsAffected <= 0 || r.Error != nil {
		transAction.Rollback()
		return false,r.Error
	}
	transAction.Commit()
	return true,nil
}

func (self *UserCashAccountService) UpdateByUserId(userCashAccount *mngModel.UserCashAccount) (bool, error) {
	transAction := common.SodaMngDB_WR.Begin()
	r := transAction.Model(&mngModel.UserCashAccount{}).Where("user_id = ?", userCashAccount.UserId).Updates(userCashAccount)
	if r.Error != nil {
		transAction.Rollback()
		return false, r.Error
	}
	//对有可能为0的值进行单独更新
	var value_zero = make(map[string]interface{})
	value_zero["type"] = userCashAccount.Type
	value_zero["province_id"] = userCashAccount.ProvinceId
	value_zero["city_id"] = userCashAccount.CityId
	value_zero["mobile"] = userCashAccount.Mobile

	//再单独更新一次type避免为0时更新不了
	r = transAction.Model(&mngModel.UserCashAccount{}).Where("user_id = ?", userCashAccount.UserId).Updates(value_zero)
	if r.Error != nil {
		transAction.Rollback()
		return false, r.Error
	}
	transAction.Commit()
	return true, r.Error
}

func (self *UserCashAccountService) ListByUserIds(userIds string) (*[]*mngModel.UserCashAccount, error) {
	list := &[]*mngModel.UserCashAccount{}
	r := common.SodaMngDB_R.Model(&mngModel.UserCashAccount{}).Where("user_id in (?)", userIds).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *UserCashAccountService) List(payType ...int) (*[]*mngModel.UserCashAccount, error) {
	list := &[]*mngModel.UserCashAccount{}
	r := common.SodaMngDB_R.Model(&mngModel.UserCashAccount{}).Where("type in (?)", payType).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}
