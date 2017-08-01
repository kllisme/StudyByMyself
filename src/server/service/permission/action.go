package permission

import (
	"maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
)

type ActionService struct {

}

func (self *ActionService)GetByID(id int) (*permission.Action, error) {
	action := permission.Action{}
	err := common.SodaMngDB_R.Where(id).Find(&action).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (self *ActionService)Query(action *permission.Action) (*[]*permission.Action, error) {
	actionList := make([]*permission.Action, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if action.HandlerName != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("handler_name like (?)", "%" + action.HandlerName + "%")
		})
	}
	if action.Method != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("method = ?", action.Method)
		})
	}
	if err := db.Scopes(scopes...).Find(&actionList).Error; err != nil {
		return nil, err
	}
	return &actionList, nil

}

func (self *ActionService)GetListByIDs(ids ...interface{}) (*[]*permission.Action, error) {
	actionList := make([]*permission.Action, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Find(&actionList).Error
	if err != nil {
		return nil, err
	}
	return &actionList, nil
}

func (self *ActionService)Create(action *permission.Action) error {
	err := common.SodaMngDB_WR.Create(action).Error
	return err
}

func (self *ActionService)Delete(id int) error {
	err := common.SodaMngDB_WR.Delete(&permission.Action{}, id).Error
	return err
}

func (self *ActionService)Update(action *permission.Action) (*permission.Action, error) {
	if err := common.SodaMngDB_WR.Save(&action).Error; err != nil {
		return nil, err
	}
	return action, nil
}
