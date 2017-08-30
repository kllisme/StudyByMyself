package public

import (
	"maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
)

type ApplicationService struct {

}

func (self *ApplicationService)GetByID(id int) (*public.Application, error) {
	application := public.Application{}
	err := common.SodaMngDB_R.Where(id).Find(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (self *ApplicationService)Paging(page int, perPage int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	applicationList := make([]*public.Application, 0)
	if err := common.SodaMngDB_R.Model(&public.Application{}).Count(&pagination.Pagination.Total).Offset((page - 1) * perPage).Limit(perPage).Order("id desc").Find(&applicationList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = (page - 1) * perPage + 1
	pagination.Pagination.To = perPage * page
	if pagination.Pagination.To > pagination.Pagination.Total {
		pagination.Pagination.To = pagination.Pagination.Total
	}
	pagination.Objects = applicationList
	return &pagination, nil

}

func (self *ApplicationService)Create(application *public.Application) (*public.Application, error) {
	err := common.SodaMngDB_WR.Create(&application).Error
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (self *ApplicationService)Update(application *public.Application) (*public.Application, error) {
	_application := map[string]interface{}{
		"name":application.Name,
		"description":application.Description,
	}
	if err := common.SodaMngDB_WR.Model(&public.Application{}).Where(application.ID).Updates(_application).Scan(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (self *ApplicationService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Delete(&public.Application{}, id).Error; err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Where("app_id = ?", id).Delete(&public.ADSpace{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

