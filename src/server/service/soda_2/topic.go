package soda_2

import (
	twoModel "maizuo.com/soda/erp/api/src/server/model/soda_2"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/entity"
	"maizuo.com/soda/erp/api/src/server/payload"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type TopicService struct {

}

func (self *TopicService)GetByID(id int) (*twoModel.Topic, error) {
	topic := twoModel.Topic{}
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
	topicList := make([]*twoModel.Topic, 0)
	channelList := make([]*twoModel.Channel, 0)
	db := common.Soda2DB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	if keywords != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("title like (?) or content like (?)", "%" + keywords + "%", "%" + keywords + "%")
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
	if err := db.Model(&twoModel.Topic{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&topicList).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&twoModel.Channel{}).Find(&channelList).Error; err != nil {
		return nil, err
	}
	for _, topic := range topicList {
		for _, channel := range channelList {
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
		province := mngModel.Province{}
		err := common.SodaMngDB_R.Where("code = ?", provinceID).Find(&province).Error
		if err != nil {
			return nil, err
		}
		cityIDs := make([]int, 0)

		err = common.SodaMngDB_R.Model(&mngModel.City{}).Where("parent_code = ?", provinceID).Pluck("code", &cityIDs).Error
		if err != nil {
			return nil, err
		}
		pagination.Pagination.Total= int(db.Raw("SELECT city_id,city_name,count(distinct school_name) as school_count,count(*) as topic_count FROM `topic` where city_id in (?) GROUP BY city_id", cityIDs).Scan(&circleList).RowsAffected)

		if err := db.Raw("SELECT city_id,city_name,count(distinct school_name) as school_count,count(*) as topic_count FROM `topic` where city_id in (?) GROUP BY city_id ORDER BY topic_count desc LIMIT ? OFFSET ?", cityIDs, limit, offset).Scan(&circleList).Error; err != nil {
			return nil, err
		}

	} else {
		pagination.Pagination.Total= int(db.Raw("SELECT city_id,city_name,count(distinct school_name) as school_count,count(*) as topic_count FROM `topic` GROUP BY city_id ").Scan(&circleList).RowsAffected)

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
		schoolIDs := make([]int, 0)
		err := common.SodaMngDB_R.Model(&mngModel.School{}).Where("city_code = ?", circle.CityID).Pluck("id", &schoolIDs).Error
		if err != nil {
			return nil, err
		}
		common.Soda2DB_R.Table("user").Where("school_id in (?)", schoolIDs).Count(&circle.UserCount)
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
	result := make([]*int, 0)
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

func (self *TopicService)UpdateChannel(entity *twoModel.Topic) (*twoModel.Topic, error) {
	_topic := map[string]interface{}{
		"channel_id":entity.ChannelID,
		"channel_title":entity.ChannelTitle,
	}
	if err := common.Soda2DB_WR.Model(&twoModel.Topic{}).Where(entity.ID).Updates(_topic).Scan(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (self *TopicService)UpdateStatus(entity *twoModel.Topic) (*twoModel.Topic, error) {
	_topic := map[string]interface{}{
		"status":entity.Status,
	}
	if err := common.Soda2DB_WR.Model(&twoModel.Topic{}).Where(entity.ID).Updates(_topic).Scan(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
