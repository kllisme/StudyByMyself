package soda_manager

import (
	mngService "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/entity"
)

type ActionService struct {

}

func (self *ActionService)GetByID(id int) (*mngService.Action, error) {
	action := mngService.Action{}
	err := common.SodaMngDB_R.Where(id).Find(&action).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (self *ActionService)Paging(offset int, limit int, handlerName string, method string) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	actionList := make([]*mngService.Action, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if handlerName != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("handler_name like (?)", "%" + handlerName + "%")
		})
	}
	if method != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("method like (?)","%" + method + "%")
		})
	}
	if err := db.Model(&mngService.Action{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&actionList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = actionList
	return &pagination, nil

}

func (self *ActionService)GetListByIDs(ids ...interface{}) (*[]*mngService.Action, error) {
	actionList := make([]*mngService.Action, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Order("id desc").Find(&actionList).Error
	if err != nil {
		return nil, err
	}
	return &actionList, nil
}

func (self *ActionService)Create(action *mngService.Action) error {
	err := common.SodaMngDB_WR.Create(action).Error
	return err
}

func (self *ActionService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&mngService.Action{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Unscoped().Where("action_id = ?", id).Delete(&mngService.PermissionActionRel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (self *ActionService)Update(action *mngService.Action) (*mngService.Action, error) {
	_action := map[string]interface{}{
		"api":action.API,
		"handler_name":action.HandlerName,
		"description":action.Description,
		"method":action.Method,
	}
	if err := common.SodaMngDB_WR.Model(&mngService.Action{}).Where(action.ID).Updates(_action).Scan(action).Error; err != nil {
		return nil, err
	}
	return action, nil
}
