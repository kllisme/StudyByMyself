package permission

import (
	"maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/jinzhu/gorm"
)

type ElementService struct {

}

func (self *ElementService)GetListByIDs(ids ...interface{}) (*[]*permission.Element, error) {
	elementList := make([]*permission.Element, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Find(&elementList).Error
	if err != nil {
		return nil, err
	}
	return &elementList, nil
}

func (self *ElementService)GetByID(id int) (*permission.Element, error) {
	element := permission.Element{}
	err := common.SodaMngDB_R.Where(id).Find(&element).Error
	if err != nil {
		return nil, err
	}
	return &element, nil
}

func (self *ElementService)Paging(page int, perPage int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	elementList := make([]*permission.Element, 0)
	db := common.SodaMngDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if err := db.Model(&permission.Element{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset((page - 1) * perPage).Limit(perPage).Find(&elementList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = (page - 1) * perPage
	pagination.Pagination.To = perPage * page
	if pagination.Pagination.To > pagination.Pagination.Total {
		pagination.Pagination.To = pagination.Pagination.Total
	}
	pagination.Objects = elementList
	return &pagination, nil

}

func (self *ElementService)Create(element *permission.Element) (*permission.Element, error) {
	err := common.SodaMngDB_WR.Create(&element).Error
	if err != nil {
		return nil, err
	}
	return element, nil
}

func (self *ElementService)Update(element *permission.Element) (*permission.Element, error) {
	_element := map[string]interface{}{
		"name":element.Name,
		"reference":element.Reference,
	}
	if err := common.SodaMngDB_WR.Model(&permission.Element{}).Where(element.ID).Updates(_element).Scan(element).Error; err != nil {
		return nil, err
	}
	return element, nil
}

func (self *ElementService)Delete(id int) error {
	err := common.SodaMngDB_WR.Delete(&permission.Element{}, id).Error
	return err
}
