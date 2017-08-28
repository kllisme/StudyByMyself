package community

import (
	"maizuo.com/soda/erp/api/src/server/model/community"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/entity"
	"maizuo.com/soda/erp/api/src/server/payload"
	"maizuo.com/soda/erp/api/src/server/model/public"
)

type TopicService struct {

}

func (self *TopicService)GetByID(id int) (*community.Topic, error) {
	topic := community.Topic{}
	err := common.SodaDB_R.Where(id).Find(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}
//
//func (self *TopicService)GetAll() (*[]*community.Topic, error) {
//	topicList := make([]*community.Topic, 0)
//	if err := common.SodaDB_R.Order("id desc").Find(&topicList).Error; err != nil {
//		return nil, err
//	}
//	return &topicList, nil
//
//}

func (self *TopicService)Paging(cityID int, keywords string, schoolName string, channelID int, status int, page int, perPage int, userIDs []int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	topicList := make([]*community.Topic, 0)
	db := common.SodaDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if keywords != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("title like (?)", "%" + keywords + "%").Or("content like (?)", "%" + keywords + "%")
		})
	}
	if schoolName != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("school_name like (?)", "%" + schoolName + "%")
		})
	}
	if status != -1 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", status)
		})
	}
	if cityID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("city_id = ?", cityID)
		})
	}
	if channelID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("channel_id = ?", channelID)
		})
	}
	common.Logger.Debugf("%#v", userIDs)
	if len(userIDs) != 0 {
		common.Logger.Debugf("------------------------------")
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id in (?)", userIDs)
		})
	}
	if err := db.Model(&community.Topic{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset((page - 1) * perPage).Limit(perPage).Order("id desc").Find(&topicList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = (page - 1) * perPage + 1
	pagination.Pagination.To = perPage * page
	if pagination.Pagination.To > pagination.Pagination.Total {
		pagination.Pagination.To = pagination.Pagination.Total
	}
	pagination.Objects = topicList
	return &pagination, nil
}

func (self *TopicService)PagingCircle(page int, perPage int, provinceID int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}

	circleList := make([]*payload.Circle, 0)
	db := common.SodaDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if provinceID != 0 {
		region:=public.Region{}
		err := common.SodaMngDB_R.Where("id = ?", provinceID).Find(&region).Error
		if err != nil {
			return nil, err
		}
		cityIDs := make([]int,0)
		if region.LevelName == "市" {
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where("city_id = ?", provinceID)
			})
		} else {
			err := common.SodaMngDB_R.Table("region").Where("parent_id = ? and level = 2", provinceID).Pluck("id",&cityIDs).Error
			if err != nil {
				return nil, err
			}
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where("city_id in (?)", cityIDs)
			})
		}

	}

	if err := db.Table("2_topic").Select("city_id,city_name,count(distinct school_name) as school_count,count(*) as topic_count").Scopes(scopes...).Group("city_id").Order("topic_count desc").Offset((page - 1) * perPage).Limit(perPage).Find(&circleList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = (page - 1) * perPage + 1
	pagination.Pagination.To = perPage * page
	if pagination.Pagination.To > pagination.Pagination.Total {
		pagination.Pagination.To = pagination.Pagination.Total
	}
	//TODO 完善学校所在圈子的逻辑
	for index, circle := range circleList {
		common.SodaDB_R.Table("2_user").Where("school_id in (?)", circle.CityID).Count(&circle.UserCount)
		circle.Order = pagination.Pagination.From + index
	}

	pagination.Objects = circleList
	return &pagination, nil
}

//func (self *TopicService)GetListByIDs(ids ...interface{}) (*[]*community.Topic, error) {
//	topicList := make([]*community.Topic, 0)
//	err := common.SodaDB_R.Where("id in (?)", ids...).Order("id desc").Find(&topicList).Error
//	if err != nil {
//		return nil, err
//	}
//	return &topicList, nil
//}

//func (self *TopicService)Create(topic *community.Topic) (*community.Topic, error) {
//	err := common.SodaMngDB_WR.Create(&topic).Error
//	if err != nil {
//		return nil, err
//	}
//	return topic, nil
//}

func (self *TopicService)CountByCityIDs(cityIDs ...interface{}) (int, error) {
	count := 0
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if  len(cityIDs) != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("city_id in (?)", cityIDs...)
		})
	}
	err := common.SodaDB_R.Table("2_topic").Scopes(scopes...).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (self *TopicService)CountCities() (int, error) {
	result := make([]*int,0)
	count := int(common.SodaDB_R.Table("2_topic").Select("distinct city_id").Scan(&result).RowsAffected)
	return count, nil
}

//func (self *TopicService)Delete(id int) error {
//	tx := common.SodaMngDB_WR.Begin()
//	if err := tx.Unscoped().Delete(&community.Topic{}, id).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	tx.Commit()
//	return nil
//}

func (self *TopicService)UpdateChannel(entity *community.Topic) (*community.Topic, error) {
	_topic := map[string]interface{}{
		"channel_id":entity.ChannelID,
		"channel_title":entity.ChannelTitle,
	}
	if err := common.SodaDB_WR.Model(&community.Topic{}).Where(entity.ID).Updates(_topic).Scan(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (self *TopicService)UpdateStatus(entity *community.Topic) (*community.Topic, error) {
	_topic := map[string]interface{}{
		"status":entity.Status,
	}
	if err := common.SodaDB_WR.Model(&community.Topic{}).Where(entity.ID).Updates(_topic).Scan(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
