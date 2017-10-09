package soda

import (
	//"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/common"
	sodaModel "maizuo.com/soda/erp/api/src/server/model/soda"
)

type UserService struct {
}

//校验所有信息
func (self *UserService) CheckInfo(user *sodaModel.User) (*sodaModel.User, error) {
	db := common.SodaDB_R
	result := sodaModel.User{}
	if err := db.Where(user).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

//func (self *UserService) GetById(id int) (*sodaModel.User, error) {
//	db := common.SodaDB_R
//	user := sodaModel.User{}
//	if err := db.Where(id).First(&user).Error; err != nil {
//		return nil, err
//	}
//	return &user, nil
//}

func (self *UserService) GetByMobile(mobile string) (*sodaModel.User, error) {
	db := common.SodaDB_R
	result := sodaModel.User{}
	if err := db.Unscoped().Where(&sodaModel.User{Mobile:mobile}).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

//func (self *UserService) GetList(name string, account string, ids []int, roleID int) (*[]*sodaModel.User, error) {
//	userList := make([]*sodaModel.User, 0)
//	db := common.SodaDB_R
//	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
//	if name != "" {
//		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
//			return db.Where("name = ?", name)
//		})
//	}
//	if account != "" {
//		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
//			return db.Where("account = ?", account)
//		})
//	}
//	if len(ids) != 0 {
//		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
//			return db.Where("id in (?)", ids)
//		})
//	}
//	if roleID != 0 {
//		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
//			return db.Joins("left join erp_user_role_rel rel on rel.user_id = user.id and rel.role_id = ?", roleID)
//		})
//	}
//	if err := db.Model(&sodaModel.User{}).Scopes(scopes...).Order("id desc").Find(&userList).Error; err != nil {
//		return nil, err
//	}
//	return &userList, nil
//}

//func (self *UserService) Update(user *sodaModel.User) (*sodaModel.User, error) {
//	_user := map[string]interface{}{
//		"name":      user.Name,
//		"contact":   user.Contact,
//		"mobile":    user.Mobile,
//		"telephone": user.Telephone,
//		"address":   user.Address,
//		"email":     user.Email,
//	}
//	if err := common.SodaDB_WR.Model(&sodaModel.User{}).Where("id = ?", user.ID).Updates(_user).Scan(user).Error; err != nil {
//		return nil, err
//	}
//	return user, nil
//}

func (self *UserService) ChangePassword(mobile string, password string) (*sodaModel.Account, error) {
	account := sodaModel.Account{}
	err := common.SodaDB_WR.Model(&sodaModel.Account{}).Where("mobile = ?",mobile).Updates(&sodaModel.Account{Password: password}).Scan(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, err
}
