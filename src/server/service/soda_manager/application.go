package soda_manager

import (
	mngService "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
)

type ApplicationService struct {

}

func (self *ApplicationService)GetByID(id int) (*mngService.Application, error) {
	application := mngService.Application{}
	err := common.SodaMngDB_R.Where(id).Find(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (self *ApplicationService)Paging(offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	applicationList := make([]*mngService.Application, 0)
	if err := common.SodaMngDB_R.Model(&mngService.Application{}).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&applicationList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = applicationList
	return &pagination, nil

}

func (self *ApplicationService)Create(application *mngService.Application) (*mngService.Application, error) {
	err := common.SodaMngDB_WR.Create(&application).Error
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (self *ApplicationService)Update(application *mngService.Application) (*mngService.Application, error) {
	_application := map[string]interface{}{
		"name":application.Name,
		"description":application.Description,
	}
	if err := common.SodaMngDB_WR.Model(&mngService.Application{}).Where(application.ID).Updates(_application).Scan(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (self *ApplicationService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Delete(&mngService.Application{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Where("app_id = ?", id).Delete(&mngService.ADPosition{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

