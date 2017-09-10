package public

import (
	"maizuo.com/soda/erp/api/src/server/common"
	model "maizuo.com/soda/erp/api/src/server/model/public"
)

type SchoolService struct {

}

func (self *SchoolService)GetListByIDs(ids ...interface{}) (*[]*model.School, error) {
	schoolList := make([]*model.School, 0)
	err := common.SodaMngDB_R.Where("id in (?)", ids...).Order("id desc").Find(&schoolList).Error
	if err != nil {
		return nil, err
	}
	return &schoolList, nil
}

func (self *SchoolService)GetByID(id int) (*model.School, error) {
	school := model.School{}
	err := common.SodaMngDB_R.Where(id).Find(&school).Error
	if err != nil {
		return nil, err
	}
	return &school, nil
}


func (self *SchoolService)Create(school *model.School) (*model.School, error) {
	err := common.SodaMngDB_WR.Create(&school).Error
	if err != nil {
		return nil, err
	}
	return school, nil
}

func (self *SchoolService)Update(school *model.School) (*model.School, error) {
	if err := common.SodaMngDB_WR.Model(&model.School{}).Save(&school).Where(school.ID).Error; err != nil {
		return nil, err
	}
	return school, nil
}

func (self *SchoolService)Delete(id int) error {
	tx := common.SodaMngDB_WR.Begin()
	if err := tx.Unscoped().Delete(&model.School{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

