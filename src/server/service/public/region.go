package public

import (
	"maizuo.com/soda/erp/api/src/server/common"
	model "maizuo.com/soda/erp/api/src/server/model/public"
)

type RegionService struct {

}


func (self *RegionService)GetByID(id int) (*model.Region, error) {
	region := model.Region{}
	err := common.SodaMngDB_R.Where(id).Find(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (self *RegionService)GetProvinces() (*[]*model.Region, error) {
	regionList := make([]*model.Region,0)
	err := common.SodaMngDB_R.Where("level = 1 ").Find(&regionList).Error
	if err != nil {
		return nil, err
	}
	return &regionList, nil
}

func (self *RegionService)GetCities(provinceID int) (*[]*model.Region, error) {
	regionList := make([]*model.Region,0)
	region:=model.Region{}
	err := common.SodaMngDB_R.Where("parent_id = ? and level = 2", provinceID).Find(&region).Error
	if err != nil {
		return nil, err
	}
	if region.LevelName == "市" {
		regionList = append(regionList, &region)
	} else {
		err := common.SodaMngDB_R.Where("parent_id = ? and level = 2", provinceID).Find(&regionList).Error
		if err != nil {
			return nil, err
		}
	}
	return &regionList, nil
}

func (self *RegionService)GetRegions(parentID int) (*[]*model.Region, error) {
	regionList := make([]*model.Region,0)
	err := common.SodaMngDB_R.Where("parent_id = ? and level = 3", parentID).Find(&regionList).Error
	if err != nil {
		return nil, err
	}
	return &regionList, nil
}
