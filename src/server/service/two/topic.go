package two

import (
	"maizuo.com/soda/erp/api/src/server/model/two"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/entity"
	"maizuo.com/soda/erp/api/src/server/payload"
	"maizuo.com/soda/erp/api/src/server/model/public"
)

type TopicService struct {

}

func (self *TopicService)GetByID(id int) (*two.Topic, error) {
	topic := two.Topic{}
	err := common.Soda2DB_R.Where(id).Find(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}
//
//func (self *TopicService)GetAll() (*[]*two.Topic, error) {
//	topicList := make([]*two.Topic, 0)
//	if err := common.SodaDB_R.Order("id desc").Find(&topicList).Error; err != nil {
//		return nil, err
//	}
//	return &topicList, nil
//
//}

func (self *TopicService)Paging(cityID int, keywords string, schoolName string, channelID int, status int, offset int, limit int, userIDs []int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	topicList := make([]*two.Topic, 0)
	channelList := make([]*two.Channel, 0)
	db := common.Soda2DB_R
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
	if len(userIDs) != 0 {
		common.Logger.Debugf("------------------------------")
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id in (?)", userIDs)
		})
	}
	if err := db.Model(&two.Topic{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&topicList).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&two.Channel{}).Find(&channelList).Error; err != nil {
		return nil, err
	}
	for _,topic := range topicList {
		for _,channel := range channelList {
			if topic.ChannelID == channel.ID {
				topic.ChannelTitle = channel.Title
			}
		}
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = topicList
	return &pagination, nil
}

func (self *TopicService)PagingCircle(offset int, limit int, provinceID int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}

	circleList := make([]*payload.Circle, 0)
	db := common.Soda2DB_R
	if provinceID != 0 {
		region := public.Region{}
		err := common.SodaMngDB_R.Where("id = ?", provinceID).Find(&region).Error
		if err != nil {
			return nil, err
		}
		cityIDs := make([]int, 0)
		if region.LevelName == "市" {
			if err := db.Raw("SELECT city_id,city_name,count(distinct school_name) as school_count,count(*) as topic_count FROM `topic` where city_id = ? GROUP BY city_id ORDER BY topic_count desc LIMIT ? OFFSET ?", provinceID, limit, offset).Scan(&circleList).Error; err != nil {
				return nil, err
			}
		} else {
			err := common.SodaMngDB_R.Table("region").Where("parent_id = ? and level = 2", provinceID).Pluck("id", &cityIDs).Error
			if err != nil {
				return nil, err
			}
			if err := db.Raw("SELECT city_id,city_name,count(distinct school_name) as school_count,count(*) as topic_count FROM `topic` where city_id in (?) GROUP BY city_id ORDER BY topic_count desc LIMIT ? OFFSET ?", cityIDs, limit, offset).Scan(&circleList).Error; err != nil {
				return nil, err
			}

		}

	} else {
		if err := db.Raw("SELECT city_id,city_name,count(distinct school_name) as school_count,count(*) as topic_count FROM `topic` GROUP BY city_id ORDER BY topic_count desc LIMIT ? OFFSET ?", limit, offset).Scan(&circleList).Error; err != nil {
			return nil, err
		}
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	//TODO 完善学校所在圈子的逻辑
	for index, circle := range circleList {
		common.Soda2DB_R.Table("user").Where("school_id in (?)", circle.CityID).Count(&circle.UserCount)
		circle.Order = pagination.Pagination.From + index
	}

	pagination.Objects = circleList
	return &pagination, nil
}

//func (self *TopicService)GetListByIDs(ids ...interface{}) (*[]*two.Topic, error) {
//	topicList := make([]*two.Topic, 0)
//	err := common.SodaDB_R.Where("id in (?)", ids...).Order("id desc").Find(&topicList).Error
//	if err != nil {
//		return nil, err
//	}
//	return &topicList, nil
//}

//func (self *TopicService)Create(topic *two.Topic) (*two.Topic, error) {
//	err := common.SodaMngDB_WR.Create(&topic).Error
//	if err != nil {
//		return nil, err
//	}
//	return topic, nil
//}

func (self *TopicService)CountByCityIDs(cityIDs ...interface{}) (int, error) {
	count := 0
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if len(cityIDs) != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("city_id in (?)", cityIDs...)
		})
	}
	err := common.Soda2DB_R.Table("topic").Scopes(scopes...).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (self *TopicService)CountCities() (int, error) {
	result := make([]*int,0)
	count := int(common.Soda2DB_R.Table("topic").Select("distinct city_id").Scan(&result).RowsAffected)
	return count, nil
}

//func (self *TopicService)Delete(id int) error {
//	tx := common.SodaMngDB_WR.Begin()
//	if err := tx.Unscoped().Delete(&two.Topic{}, id).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	tx.Commit()
//	return nil
//}

func (self *TopicService)UpdateChannel(entity *two.Topic) (*two.Topic, error) {
	_topic := map[string]interface{}{
		"channel_id":entity.ChannelID,
		"channel_title":entity.ChannelTitle,
	}
	if err := common.Soda2DB_WR.Model(&two.Topic{}).Where(entity.ID).Updates(_topic).Scan(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (self *TopicService)UpdateStatus(entity *two.Topic) (*two.Topic, error) {
	_topic := map[string]interface{}{
		"status":entity.Status,
	}
	if err := common.Soda2DB_WR.Model(&two.Topic{}).Where(entity.ID).Updates(_topic).Scan(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
